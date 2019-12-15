// Sample
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/m-lab/etl/active"

	"cloud.google.com/go/storage"
	"github.com/GoogleCloudPlatform/google-cloud-go-testing/storage/stiface"
	"google.golang.org/api/option"

	"github.com/m-lab/etl/etl"
	"github.com/m-lab/etl/metrics"
	"github.com/m-lab/etl/worker"
	"github.com/m-lab/go/prometheusx"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	// Enable profiling. For more background and usage information, see:
	//   https://blog.golang.org/profiling-go-programs
	_ "net/http/pprof"

	// Enable exported debug vars.  See https://golang.org/pkg/expvar/
	_ "expvar"
)

func init() {
	// Always prepend the filename and line number.
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

// Task Queue can always submit to an admin restricted URL.
//   login: admin
// Return 200 status code.
// Track reqeusts that last longer than 24 hrs.
// Is task handling idempotent?

// Useful headers added by AppEngine when sending Tasks via Push.
//   X-AppEngine-QueueName
//   X-AppEngine-TaskETA
//   X-AppEngine-TaskName
//   X-AppEngine-TaskRetryCount
//   X-AppEngine-TaskExecutionCount

// TODO(gfr) Add either a black list or a white list for the environment
// variables, so we can hide sensitive vars. https://github.com/m-lab/etl/issues/384
func Status(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<html><body>\n")
	fmt.Fprintf(w, "<p>NOTE: This is just one of potentially many instances.</p>\n")
	commit := os.Getenv("COMMIT_HASH")
	if len(commit) >= 8 {
		fmt.Fprintf(w, "Release: %s <br>  Commit: <a href=\"https://github.com/m-lab/etl/tree/%s\">%s</a><br>\n",
			os.Getenv("RELEASE_TAG"), os.Getenv("COMMIT_HASH"), os.Getenv("COMMIT_HASH")[0:7])
	} else {
		fmt.Fprintf(w, "Release: %s   Commit: unknown\n", os.Getenv("RELEASE_TAG"))
	}

	fmt.Fprintf(w, "<p>Workers: %d / %d</p>\n", atomic.LoadInt32(&inFlight), maxInFlight)
	env := os.Environ()
	for i := range env {
		fmt.Fprintf(w, "%s</br>\n", env[i])
	}
	fmt.Fprintf(w, "</body></html>\n")
}

// Basic throttling to restrict the number of tasks in flight.
const defaultMaxInFlight = 20

// This limits the number of workers available for externally requested single task files.
var maxInFlight int32 // Max number of concurrent workers (and tasks in flight).
var inFlight int32    // Current number of tasks in flight.

// Returns true if request should be rejected.
// If the max concurrency (MC) exceeds (or matches) the instances*workers, then
// most requests will be rejected, until the median number of workers is
// less than the throttle.
// ** So we should set max instances (MI) * max workers (MW) > max concurrency.
//
// We also want max_concurrency high enough that most instances have several
// jobs.  With MI=20, MW=25, MC=100, the average workers/instance is only 4, and
// we end up with many instances starved, so AppEngine was removing instances even
// though the queue throughput was poor.
// ** So we probably want MC/MI > MW/2, to prevent starvation.
//
// For now, assuming:
//    MC: 180,  MI: 20, MW: 10
//
// TODO - replace the atomic with a channel based semaphore and non-blocking
// select.
func shouldThrottle() bool {
	if atomic.AddInt32(&inFlight, 1) > maxInFlight {
		atomic.AddInt32(&inFlight, -1)
		return true
	}
	return false
}

func decrementInFlight() {
	atomic.AddInt32(&inFlight, -1)
}

// TODO(gfr) unify counting for http and pubsub paths?
func handleRequest(rwr http.ResponseWriter, rq *http.Request) {
	// This will add metric count and log message from any panic.
	// The panic will still propagate, and http will report it.
	defer func() {
		metrics.CountPanics(recover(), "handleRequest")
	}()

	// Throttle by grabbing a semaphore from channel.
	if shouldThrottle() {
		metrics.TaskCount.WithLabelValues("unknown", "worker", "TooManyRequests").Inc()
		rwr.WriteHeader(http.StatusTooManyRequests)
		fmt.Fprintf(rwr, `{"message": "Too many tasks."}`)
		return
	}

	// Decrement counter when worker finishes.
	defer decrementInFlight()

	var err error
	retryCountStr := rq.Header.Get("X-AppEngine-TaskRetryCount")
	retryCount := 0
	if retryCountStr != "" {
		retryCount, err = strconv.Atoi(retryCountStr)
		if err != nil {
			log.Printf("Invalid retries string: %s\n", retryCountStr)
		}
	}
	executionCountStr := rq.Header.Get("X-AppEngine-TaskExecutionCount")
	executionCount := 0
	if executionCountStr != "" {
		executionCount, err = strconv.Atoi(executionCountStr)
		if err != nil {
			log.Printf("Invalid execution count string: %s\n", executionCountStr)
		}
	}
	etaUnixStr := rq.Header.Get("X-AppEngine-TaskETA")
	etaUnixSeconds := float64(0)
	if etaUnixStr != "" {
		etaUnixSeconds, err = strconv.ParseFloat(etaUnixStr, 64)
		if err != nil {
			log.Printf("Invalid eta string: %s\n", etaUnixStr)
		}
	}
	etaTime := time.Unix(int64(etaUnixSeconds), 0) // second granularity is sufficient.
	age := time.Since(etaTime)

	rq.ParseForm()
	// Log request data.
	for key, value := range rq.Form {
		log.Printf("Form:   %q == %q\n", key, value)
	}

	rawFileName := rq.FormValue("filename")
	status, msg := subworker(rawFileName, executionCount, retryCount, age)
	rwr.WriteHeader(status)
	fmt.Fprintf(rwr, msg)
}

func subworker(rawFileName string, executionCount, retryCount int, age time.Duration) (status int, msg string) {
	// TODO(dev) Check how many times a request has already been attempted.

	var err error
	// This handles base64 encoding, and requires a gs:// prefix.
	fn, err := etl.GetFilename(rawFileName)
	if err != nil {
		metrics.TaskCount.WithLabelValues("unknown", "worker", "BadRequest").Inc()
		log.Printf("Invalid filename: %s\n", fn)
		return http.StatusBadRequest, `{"message": "Invalid filename."}`
	}

	// TODO(dev): log the originating task queue name from headers.
	log.Printf("Received filename: %q  Retries: %d, Executions: %d, Age: %5.2f hours\n",
		fn, retryCount, executionCount, age.Hours())

	status, err = worker.ProcessTask(fn)
	if err == nil {
		msg = `{"message": "Success"}`
	} else {
		msg = fmt.Sprintf(`{"message": "%s"}`, err.Error())
	}
	return
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// TODO(soltesz): provide a real health check.
	fmt.Fprint(w, "ok")
}

func setMaxInFlight() {
	maxInFlightString, ok := os.LookupEnv("MAX_WORKERS")
	if ok {
		maxInFlightInt, err := strconv.Atoi(maxInFlightString)
		if err == nil {
			maxInFlight = int32(maxInFlightInt)
		} else {
			log.Println("MAX_WORKERS not configured.  Using 20.")
			maxInFlight = defaultMaxInFlight
		}
	} else {
		log.Println("MAX_WORKERS not configured.  Using 20.")
		maxInFlight = defaultMaxInFlight
	}
}

func runFunc(o *storage.ObjectAttrs) active.Runnable {
	path := "gs://" + o.Bucket + "/" + o.Name
	return func() error {
		log.Println(path)
		_, err := worker.ProcessTask(path)

		return err
	}
}

// This is a hack, and should not generally be used.
// It has no admission control.
func handleActiveRequest(rwr http.ResponseWriter, rq *http.Request) {
	// This will add metric count and log message from any panic.
	// The panic will still propagate, and http will report it.
	defer func() {
		metrics.CountPanics(recover(), "handleActiveRequest")
	}()

	path := rq.FormValue("path")
	if len(path) == 0 {
		// TODO add metric
		rwr.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(rwr, `{"message": "Missing path"}`)
		return
	}

	client, err := storage.NewClient(context.Background(), option.WithScopes(storage.ScopeReadOnly))
	if err != nil {
		// TODO add metric
		rwr.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(rwr, "Error creating storage client")
		return
	}
	lister := active.FileListerFunc(stiface.AdaptClient(client), path)
	ctx, cancel := context.WithTimeout(context.Background(), 24*time.Hour)
	fileSource, err := active.NewGCSSource(ctx, lister, runFunc)
	if err != nil {
		cancel()
		rwr.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(rwr, fmt.Sprintf(`{"message": "Invalid path: %s %s"}`, path, err.Error()))
		return
	}

	throttle := active.NewWSTokenSource(60)

	// Run all tasks, and log error on completion.
	go func() {
		err := active.RunAll(ctx, active.Throttle(fileSource, throttle))
		cancel()

		if err != nil {
			// TODO add metric
			log.Println(path, "Had errors:", err)
		}
	}()

	rwr.WriteHeader(http.StatusOK)
	fmt.Fprintf(rwr, fmt.Sprintf(`{"message": "Processing %s"}`, path))
}

func main() {
	// Expose prometheus and pprof metrics on a separate port.
	prometheusx.MustStartPrometheus(":9090")

	http.HandleFunc("/", Status)
	http.HandleFunc("/status", Status)
	http.HandleFunc("/worker", metrics.DurationHandler("generic", handleRequest))
	http.HandleFunc("/active", metrics.DurationHandler("generic", handleActiveRequest))
	http.HandleFunc("/_ah/health", healthCheckHandler)

	// Enable block profiling
	runtime.SetBlockProfileRate(1000000) // One event per msec.

	setMaxInFlight()

	// We also setup another prometheus handler on a non-standard path. This
	// path name will be accessible through the AppEngine service address,
	// however it will be served by a random instance.
	http.Handle("/random-metrics", promhttp.Handler())
	http.ListenAndServe(":8080", nil)
}
