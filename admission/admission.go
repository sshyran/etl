// Package admission regulates the admission of new work items, based on cpu utilization.
package admission

// Initially we just have free tokens = ...

import (
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/shirou/gopsutil/cpu"
)

// Store past 10 sample
var stats = make([]cpu.TimesStat, 0, 10)

var lastBusy int32

var lock sync.Mutex
var inFlight int32
var available int32

// GetToken gets a token if one is available.
func GetToken() bool {
	return false
}

// ReturnToken releases a token
func ReturnToken() {
}

func addToken() {

}

func removeToken() {

}

func setBusy(busy float64) {
	atomic.StoreInt32(&lastBusy, int32(busy*100))
}

func GetBusy() float64 {
	return float64(atomic.LoadInt32(&lastBusy)) / 100
}

func getAllBusy(t cpu.TimesStat) (float64, float64) {
	// IIUC linux correctly, this may double count some time.
	// Guest is included in User, and GuestNice is included in Nice.
	// Also, is Steal or Stolen counted in both the stealing thread and the stolen thread?
	//
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

func report(t1, t2 cpu.TimesStat) {
	delta := cpu.TimesStat{}
	delta.Idle = t2.Idle - t1.Idle

	delta.User = t2.User - t1.User
	delta.System = t2.System - t1.System
	delta.Nice = t2.Nice - t1.Nice
	delta.Iowait = t2.Iowait - t1.Iowait
	delta.Irq = t2.Irq - t1.Irq
	delta.Softirq = t2.Softirq - t1.Softirq

	delta.Steal = t2.Steal - t1.Steal
	delta.Stolen = t2.Stolen - t1.Stolen

	delta.Guest = t2.Guest - t1.Guest
	delta.GuestNice = t2.GuestNice - t1.GuestNice
	all, _ := getAllBusy(delta)

	log.Printf("CPU stats: %.3f/%.3f, %v\n", delta.Idle, all, delta)
}

func rotate(ts cpu.TimesStat) cpu.TimesStat {
	first := stats[0]
	start := 0
	if len(stats) > 9 {
		start = 1
	}
	stats = append(stats[start:], ts)
	return first
}

func update(ts cpu.TimesStat) cpu.TimesStat {
	first := rotate(ts)
	percent := calculateBusy(first, ts)
	setBusy(percent)
	return first
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

		ts := stat[0]
		first := update(ts)
		report(first, ts)
		<-ticker.C
	}
}
