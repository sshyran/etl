package limiter_test

import (
	"fmt"
	"log"
	"math"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/m-lab/etl/limiter"
)

func init() {
	// Always prepend the filename and line number.
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func TestBasic(t *testing.T) {
	m := limiter.StartCPUMonitor(24, 50*time.Millisecond)
	end := time.Now().Add(2 * time.Second)
	x := 5.0
	for time.Now().Before(end) {
		x = math.Sqrt(x) + 2.0
	}
	m.Kill()
}

func TestLatency(t *testing.T) {
	interval := 50 * time.Millisecond // Less than this gets really noisy.
	// The rms latency can never exceed the timer interval, so we use 1/2 the interval as the target.
	target := interval.Seconds() / 5
	log.Println("Target:", target)

	m := limiter.StartLatencyMonitor(interval, 3*interval, 10*interval)

	start := time.Now()
	end := start.Add(time.Duration(10*runtime.NumCPU()) * interval)

	var threads, threadSum, count int32
	var adds, skips int32
	once := sync.Once{}
	// Every 10 msec, decide whether to start another busy thread.
	for time.Now().Before(end) {
		th := atomic.LoadInt32(&threads)
		// Reset the average tracker when we first exceed NumCPU()
		if int(th) > runtime.NumCPU() {
			once.Do(func() {
				count = 0
				threadSum = 0
			})
		}
		count++
		threadSum += th

		r := m.Report()
		if r.Fast+.1*r.Rate < target && r.Slow < target {
			adds++
			log.Printf("%5.2f: + %2d %8.4f %8.4f %8.3f\n", time.Now().Sub(start).Seconds(), th, r.Fast, r.Slow, r.Rate)
			atomic.AddInt32(&threads, 1)
			go func() {
				// Each busy loop should run at least 3 * NumCPU * interval, so we can max out at 2 * NumCPU
				end := time.Now().Add(time.Duration(3*runtime.NumCPU()) * interval)
				x := 5.0
				for time.Now().Before(end) {
					x = math.Sqrt(x) + 2.0
				}
				th := atomic.AddInt32(&threads, -1)
				log.Printf("%5.2f: - %2d %8.4f %8.4f %8.3f\n", time.Now().Sub(start).Seconds(), th, r.Fast, r.Slow, r.Rate)
			}()
		} else {
			skips++
			log.Printf("%5.2f:   %2d %8.4f %8.4f %8.3f\n", time.Now().Sub(start).Seconds(), th, r.Fast, r.Slow, r.Rate)
		}

		time.Sleep(interval)
	}

	// We expect the algorithm to just barely fill the CPUs.
	th := atomic.LoadInt32(&threads)
	avgThreads := float32(threadSum) / float32(count)
	log.Printf("Adds: %d, Skips: %d, Threads: %d (avg = %4.1f/%d)\n", adds, skips, th, avgThreads, runtime.NumCPU())

	r := m.Report()
	if r.Fast > 2*target || r.Slow > 2*target {
		t.Error(fmt.Printf("%v", r))
	}
}
