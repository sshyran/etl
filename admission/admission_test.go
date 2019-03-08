package admission_test

import (
	"testing"
	"time"

	"github.com/m-lab/etl/admission"
)

func TestMonitor(t *testing.T) {
	admission.Monitor(time.Second)
	time.Sleep(10 * time.Second)
}
