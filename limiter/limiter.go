// package limiter implements an prioritizing admission controler.
// A client can request a Limiter, from which it can obtain tokens.
// Earlier clients are prioritized over later clients.
// Tokens are allocated based on cpu utilization.  If average utilization over past
// 10 seconds is less than A, OR over last 2 minutes is less than B, then a token
// will be granted.

package limiter

import (
	"log"
	"runtime"
	"sync"
	"syscall"
	"time"
)

var monitor = StartCPUMonitor(24, 5*time.Second)

func init() {
}

type userSys struct {
	user int64 // nanoseconds of user time
	sys  int64 // nanoseconds of system time
}

// A CPUMonitor takes CPU utilization snapshots, and reports requested interval averages.
type CPUMonitor struct {
	interval  time.Duration
	ticker    *time.Ticker
	lock      sync.Mutex
	snapshots []userSys
}

// StartCPUMonitor creates a CPUMonitor and starts it.
func StartCPUMonitor(snaps int, interval time.Duration) *CPUMonitor {
	if interval.Nanoseconds() <= 0 || snaps <= 2 {
		return nil
	}
	ticker := time.NewTicker(interval)
	snapshots := make([]userSys, snaps)
	ru := syscall.Rusage{}
	us := userSys{ru.Utime.Nano(), ru.Stime.Nano()}
	for i := range snapshots {
		snapshots[i] = us
	}

	m := CPUMonitor{ticker: ticker, snapshots: snapshots}

	go m.run()

	return &m
}

// run gathers periodic snapshots.
// Close ticker to terminate.
func (m *CPUMonitor) run() {
	for range m.ticker.C {
		ru := syscall.Rusage{}
		syscall.Getrusage(0, &ru)
		us := userSys{ru.Utime.Nano(), ru.Stime.Nano()}
		log.Println(us)
		m.lock.Lock()
		m.snapshots = append(m.snapshots[1:], us)
		m.lock.Unlock()

		a, b := m.intervals(2, 24)
		log.Printf("%5.3f %5.3f\n", a, b)
	}
}

// The average utilization for last 10 seconds (up to 5 seconds in the past), and last two minutes.
func (m *CPUMonitor) intervals(j, k int) (float64, float64) {
	m.lock.Lock()
	defer m.lock.Unlock()
	end := len(m.snapshots) - 1
	if j <= 0 || k <= 0 || j > end || k > end {
		return 0, 0
	}

	last := m.snapshots[end]
	jValues := m.snapshots[end-j]
	kValues := m.snapshots[end-k]

	lastTime := (last.user + last.sys) / 1000000      // milliseconds
	jMillis := (jValues.user + jValues.sys) / 1000000 // milliseconds
	kMillis := (kValues.user + kValues.sys) / 1000000 // milliseconds

	milliCPUPerInterval := float64(m.interval.Milliseconds() * int64(runtime.NumCPU()))
	return float64(lastTime-jMillis) / (float64(j) * milliCPUPerInterval),
		float64(lastTime-kMillis) / (float64(k) * milliCPUPerInterval)
}

var limiters = make([]*Limiter, 0, 10)

type Token struct {
	lim  Limiter // The limiter that provided the token.
	once sync.Once
}

func (t *Token) Release() {
	t.once.Do(func() {
		t.lim.release()
	})
}

type Limiter struct {
}

func (l *Limiter) Close() {

}

func (l *Limiter) release() {

}

func (l *Limiter) Get() *Token {
	return nil
}
func (l *Limiter) GetNonBlocking() *Token {
	return nil
}
