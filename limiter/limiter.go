// package limiter implements an prioritizing admission controler.
// A client can request a Limiter, from which it can obtain tokens.
// Earlier clients are prioritized over later clients.
// Tokens are allocated based on cpu utilization.  If average utilization over past
// 10 seconds is less than A, OR over last 2 minutes is less than B, then a token
// will be granted.

package limiter

import (
	"log"
	"math"
	"runtime"
	"sync"
	"syscall"
	"time"
)

// var monitor = StartCPUMonitor(24, 5*time.Second)

func init() {
}

type userSys struct {
	user int64 // nanoseconds of user time
	sys  int64 // nanoseconds of system time
}

// A CPUMonitor takes CPU utilization snapshots, and reports requested interval averages.
type CPUMonitor struct {
	size      int // Number of entries - 1.
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

	m := CPUMonitor{size: snaps - 1, interval: interval, ticker: ticker, snapshots: snapshots}

	go m.run()

	return &m
}

func (m *CPUMonitor) Kill() {
	m.ticker.Stop()
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
		log.Printf("%10.7f %10.7f\n", a, b)
	}
}

// The average utilization for last 10 seconds (up to 5 seconds in the past), and last two minutes.
func (m *CPUMonitor) intervals(j, k int) (float64, float64) {
	if j <= 0 || k <= 0 || j > m.size || k > m.size {
		return 0, 0
	}
	m.lock.Lock()

	last := m.snapshots[m.size]
	jValues := m.snapshots[m.size-j]
	kValues := m.snapshots[m.size-k]
	m.lock.Unlock()

	lastTime := float64(last.user+last.sys) / 1000000.0      // milliseconds
	jMillis := float64(jValues.user+jValues.sys) / 1000000.0 // milliseconds
	kMillis := float64(kValues.user+kValues.sys) / 1000000.0 // milliseconds

	milliCPUPerInterval := 1000.0 * (m.interval.Seconds() * float64(runtime.NumCPU()))
	log.Println(lastTime, jMillis, kMillis, milliCPUPerInterval)
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

type LatencyReport struct {
	Fast float64 // Fast filtered latency in seconds
	Slow float64 // Slow filtered latency in seconds
	Rate float64 // Rate of change of latency, in seconds/second.
}

// LatencyMonitor monitors the latency of the goroutine scheduling.
type LatencyMonitor struct {
	interval, fastInterval, slowInterval time.Duration
	ticker                               *time.Ticker
	start                                time.Time
	fastSum, slowSum                     float64 // fast and slow sumsq filters
	rate                                 float64 // rate of change, filtered at fastInterval
}

// StartLatencyMonitor creates and starts a LatencyMonitor
func StartLatencyMonitor(sample time.Duration, fast time.Duration, slow time.Duration) *LatencyMonitor {
	if sample <= 0 || fast <= sample || slow <= fast {
		return nil
	}
	m := LatencyMonitor{interval: sample, fastInterval: fast, slowInterval: slow, ticker: time.NewTicker(sample)}
	go m.run()

	return &m
}

func (m *LatencyMonitor) run() {
	last := time.Now()
	m.start = last
	for ; ; <-m.ticker.C {
		now := time.Now()

		fast := m.fastSum

		diff := now.Sub(last).Seconds() - m.interval.Seconds()
		m.fastSum = ((m.fastInterval-m.interval).Seconds()*m.fastSum + m.interval.Seconds()*(diff*diff)) / m.fastInterval.Seconds()
		m.slowSum = ((m.slowInterval-m.interval).Seconds()*m.slowSum + m.interval.Seconds()*(diff*diff)) / m.slowInterval.Seconds()
		rate := (math.Sqrt(m.fastSum) - math.Sqrt(fast)) / m.interval.Seconds()
		m.rate = ((m.fastInterval-m.interval).Seconds()*m.rate + m.interval.Seconds()*rate) / m.fastInterval.Seconds()

		last = now
	}
}

// Report prints a report and returns fast and slow metrics in seconds.
func (m *LatencyMonitor) Report() LatencyReport {
	r := LatencyReport{}
	// Ok to access concurrently, or should these be atomic?
	// Note that these are float64 values.
	r.Fast = math.Sqrt(m.fastSum)
	r.Slow = math.Sqrt(m.slowSum)
	r.Rate = m.rate
	return r
}
