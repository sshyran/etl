// Package active provides code for managing processing of an entire
// directory of task files.
package active

// The "active" task manager supervises launching of workers to process all files
// in a gs:// prefix.  It processes the files in lexographical order, and maintains
// status info in datastore.

// Design:
//  1. a token channel is passed in to ProcessAll, and used to determine how many tasks may
//     be in flight.  It is returned to the caller when there are no more tasks to start,
//     but there may still be tasks running, and tokens that will be returned later.
//  2. a doneHandler waits for task completions, and updates the state.  It starts additional
//     tasks if there are any.  When there are no more tasks, it signals ProcessAll
//     that the token channel may be returned to the caller.

// TODO:
// A. Add metrics
//
// B. Recovery and monitoring using datastore.
//
// C. Utilization based management:
//    The manager starts new tasks when either:
//   1. Two tasks have completed since the last task started.
//   2. The 10 second utilization of any single cpu falls below 80%.
//   3. The total 10 second utilization falls below 90%.

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"cloud.google.com/go/storage"
	"github.com/GoogleCloudPlatform/google-cloud-go-testing/storage/stiface"

	"github.com/m-lab/etl/cloud/gcs"
	"github.com/m-lab/etl/worker"
)

// TaskFile maintains the status of a single file.
// These are NOT thread-safe, and should only be read and modified by
// a single goroutine.
type TaskFile struct {
	// TODO - is this needed?
	path string               // Full path to object.
	obj  *storage.ObjectAttrs // Optional if completed or errored.

	failures int
	lastErr  error
}

func (tf TaskFile) String() string {
	return fmt.Sprintf("%s: %d failures, %v", tf.path, tf.failures, tf.lastErr)
}

// Path returns the full path to the file.
func (tf TaskFile) Path() string {
	return tf.path
}

// FileSource handles reading, caching, and updating a list of files,
// and tracking the processing status of each file.
type FileSource struct {
	client  stiface.Client
	project string
	prefix  string

	lock       sync.Mutex
	dispatched map[string]struct{} // To keep track of all files completed or in flight.
	pending    []*TaskFile         // Ordered list - TODO make this a channel?
	lastUpdate time.Time           // Time of last call to UpdatePending
}

// NewFileSource creates a new source for active processing.
func NewFileSource(sc stiface.Client, project string, prefix string) (*FileSource, error) {
	fs := FileSource{
		client:     sc,
		project:    project,
		prefix:     prefix,
		pending:    make([]*TaskFile, 0, 1),
		dispatched: make(map[string]struct{}, 100),
	}

	return &fs, nil
}

// AddPending puts a task file on the end of the pending list.
func (fs *FileSource) AddPending(tf *TaskFile) {
	fs.lock.Lock()
	defer fs.lock.Unlock()
	fs.pending = append(fs.pending, tf)

}

// updatePending should be called when there are no more pending tasks.
func (fs *FileSource) updatePending(ctx context.Context) error {
	fs.lock.Lock()
	defer fs.lock.Unlock()
	if len(fs.pending) > 0 {
		return nil
	}

	// Allow for a little clock skew.
	updateTime := time.Now().Add(-time.Second)
	files, _, err := gcs.GetFilesSince(context.Background(), fs.client, fs.project, fs.prefix, fs.lastUpdate)
	if err != nil {
		return err
	}
	fs.lastUpdate = updateTime

	if len(fs.pending) == 0 && cap(fs.pending) < len(files)+10 {
		fs.pending = make([]*TaskFile, 0, len(files)+10) // A few extra slots for retries.
	}
	for _, f := range files {
		if f.Prefix != "" {
			log.Println("Skipping subdirectory:", f.Prefix)
			continue // skip directories
		}
		// Append any new files that haven't already been dispatched.
		if _, exists := fs.dispatched[f.Name]; !exists {
			log.Println("Adding", "gs://"+f.Bucket+"/"+f.Name)
			tf := TaskFile{path: "gs://" + f.Bucket + "/" + f.Name, obj: f}
			fs.dispatched[f.Name] = struct{}{}
			fs.pending = append(fs.pending, &tf)
		}
	}

	return nil
}

// next returns the next pending TaskFile.  It runs Update if there
// are initially none available, and reprocesses tasks from the
// errored list if the there are still none pending.
// Caller should have already obtained a semaphore.
// Returns an error iff updatePending errored.
func (fs *FileSource) next(ctx context.Context) (*TaskFile, error) {
	err := fs.updatePending(ctx)
	if err != nil {
		return nil, err
	}
	fs.lock.Lock()
	defer fs.lock.Unlock()
	if len(fs.pending) > 0 {
		tf := fs.pending[0]
		fs.pending = fs.pending[1:]
		return tf, nil
	}
	return nil, nil
}

// ErrTaskNotFound is returned if an inflight task is not found in inFlight.
var ErrTaskNotFound = errors.New("task not found")

// TokenSource defines the interface for obtaining admission tokens
type TokenSource interface {
	Acquire(ctx context.Context) error
	Release()
}

type Dispatcher struct {
	fs *FileSource

	// The function that processes each task.
	process    func(*TaskFile) error
	retryLimit int // number of retries for failed tasks.

	// Channel to handle cleanup when a task is completed.
	done chan *TaskFile

	// Remaining fields should only be accessed while holding the lock.
	lock     sync.Mutex
	inFlight map[string]*TaskFile

	err []*TaskFile
}

// NewDispatcher creates a new Dispatcher.
func NewDispatcher(fs *FileSource, pf func(*TaskFile) error, retry int) *Dispatcher {
	return &Dispatcher{
		fs:         fs,
		retryLimit: retry,
		process:    pf,
		inFlight:   make(map[string]*TaskFile, 100),
		done:       make(chan *TaskFile, 0),
		err:        make([]*TaskFile, 0, 10),
	}
}

// updateState updates the FileSource to reflect the completion of
// a processing attempt.
// If the processing ends in an error, the task will be moved to
// the end of the pending list, unless the task has already been retried fs.retry times.
func (disp *Dispatcher) updateState(tf *TaskFile) error {
	disp.lock.Lock()
	defer disp.lock.Unlock()
	_, exists := disp.inFlight[tf.path]
	if !exists {
		log.Println("Did not find", tf.path)
		return ErrTaskNotFound
	}

	delete(disp.inFlight, tf.path)
	if tf.lastErr != nil {
		if tf.failures < disp.retryLimit {
			disp.fs.AddPending(tf)
			tf.failures++
		} else {
			disp.err = append(disp.err, tf)
		}
	}

	return tf.lastErr
}

func processTask(tf *TaskFile) error {
	_, err := worker.ProcessTask(tf.path)
	return err
}

func (disp *Dispatcher) startTask(tf *TaskFile) {
	// Add the task to inFlight map and start the task.
	disp.lock.Lock()
	defer disp.lock.Unlock()
	disp.inFlight[tf.path] = tf
	go func() {
		tf.lastErr = disp.process(tf)
		disp.done <- tf
	}()
}

func (disp *Dispatcher) startDoneHandler(tokens TokenSource) *sync.WaitGroup {
	wg := sync.WaitGroup{}
	go func() {
		for tf := range disp.done {
			log.Println("received", tf)
			err := disp.updateState(tf) // This removes the task from inFlight.
			if err != nil {
				log.Println(tf.Path(), err)
			}
			tokens.Release()
			wg.Done()
		}
	}()

	return &wg
}

// startLauncher starts the goroutines to start new tasks, and handle completions.
// It returns a sync.WaitGroup that will signal only when all jobs have completed.
func (disp *Dispatcher) startLauncher(ctx context.Context, tokens TokenSource) *sync.WaitGroup {
	wg := disp.startDoneHandler(tokens)
	wg.Add(1) // To prevent early Done() detection.
	go func() {
		for {
			// Wait for a token
			err := tokens.Acquire(ctx)
			if err != nil {
				log.Println(err)
				break // Context expired.
			}

			// If there are tokens available, start another job
			tf, err := disp.fs.next(ctx)
			if err != nil {
				log.Println(err)
			}
			if tf == nil {
				// return the token and quit
				tokens.Release()
				log.Println("No more tasks")
				break
			}

			log.Println("starting", tf)
			wg.Add(1)
			disp.startTask(tf)
		}
		wg.Done()
	}()

	return wg
}

// Errors returns a list of all TaskFile objects that ended with error.
func (disp *Dispatcher) Errors() []*TaskFile {
	return disp.err
}

// ProcessAll iterates through all the TaskFiles, processing each one.
// It may also retry any that failed the first time.
func (disp *Dispatcher) ProcessAll(ctx context.Context, fs *FileSource, tokens TokenSource) (*sync.WaitGroup, error) {
	err := fs.updatePending(ctx)
	if err != nil {
		return nil, err
	}
	// Handle tasks in parallel.
	// When the returned wg is signaled, there may still be tasks in flight, but no more
	// will be started.
	return disp.startLauncher(ctx, tokens), nil

}
