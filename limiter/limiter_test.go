package limiter_test

import (
	"log"
	"math"
	"testing"
	"time"

	"github.com/m-lab/etl/limiter"
)

func init() {
	// Always prepend the filename and line number.
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func TestBasic(t *testing.T) {
	limiter.StartCPUMonitor(24, 50*time.Millisecond)
	end := time.Now().Add(2 * time.Second)
	x := 5.0
	for time.Now().Before(end) {
		x = math.Sqrt(x) + 2.0
	}
	t.Fail()
}
