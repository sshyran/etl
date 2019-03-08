// Package admission regulates the admission of new work items, based on cpu utilization.
package admission

import (
	"log"
	"sync/atomic"
	"time"

	"github.com/shirou/gopsutil/cpu"
)

// Store past 5 minutes
var stats = make([]cpu.TimesStat, 0, 30)

var lastBusy int32

func setBusy(busy float64) {
	atomic.StoreInt32(&lastBusy, int32(busy*100))
}

func GetBusy() float64 {
	return float64(atomic.LoadInt32(&lastBusy)) / 100
}

func getAllBusy(t cpu.TimesStat) (float64, float64) {
	busy := t.User + t.System + t.Nice + t.Iowait + t.Irq +
		t.Softirq + t.Steal + t.Guest + t.GuestNice + t.Stolen
	return busy + t.Idle, busy
}

func calculateBusy(t1, t2 cpu.TimesStat) float64 {
	t1All, t1Busy := getAllBusy(t1)
	t2All, t2Busy := getAllBusy(t2)

	if t2Busy <= t1Busy {
		return 0
	}
	if t2All <= t1All {
		return 1
	}
	return (t2Busy - t1Busy) / (t2All - t1All) * 100
}

func rotate(ts cpu.TimesStat) {
	stats = append(stats, ts)
	if len(stats) > 29 {
		stats = stats[1:]
	}
}

func update(ts cpu.TimesStat) float64 {
	percent := calculateBusy(stats[0], ts)
	setBusy(percent)
	rotate(ts)
	return percent
}

func Monitor(interval time.Duration) {
	stat, err := cpu.Times(false)

	if err != nil {
		log.Println(err)
		return
	}

	rotate(stat[0])
	for {
		ticker := time.NewTicker(interval)

		stat, err := cpu.Times(false) // across all CPUs
		if err != nil {
			log.Println(err)
		}
		update(stat[0])
		log.Println("CPU percentage:", GetBusy())
		<-ticker.C
	}
}
