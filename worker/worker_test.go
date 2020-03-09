package worker_test

import (
	"log"
	"math"
	"net/http"
	"testing"

	//"github.com/m-lab/etl/metrics"
	"github.com/m-lab/etl/worker"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
)

func counterValue(m prometheus.Metric) float64 {
	var mm dto.Metric
	m.Write(&mm)
	ctr := mm.GetCounter()
	if ctr == nil {
		log.Println(mm.GetUntyped())
		return math.Inf(-1)
	}

	return *ctr.Value
}

func checkCounter(t *testing.T, c chan prometheus.Metric, expected float64) {
	m := <-c
	v := counterValue(m)
	if v != expected {
		t.Error("For", m.Desc(), "expected:", expected, "got:", v)
	}
}

func TestProcessTask(t *testing.T) {
	if testing.Short() {
		t.Log("Skipping integration test")
	}
	//filename := "gs://archive-mlab-testing/ndt/2018/05/09/20180509T101913Z-mlab1-mad03-ndt-0000.tgz"
	//filename := "gs://archive-measurement-lab/paris-traceroute/2013/05/24/20130524T000000Z-mlab3-lju01-paris-traceroute-0000.tgz"
	filename := "gs://archive-mlab-testing/host/traceroute/2019/11/15/20191115T034951.000655Z-traceroute-mlab1-tpe01-host.tgz"
	status, err := worker.ProcessTask(filename)
	if err != nil {
		t.Error(err)
	}
	if status != http.StatusOK {
		t.Error("Expected", http.StatusOK, "Got:", status)
	}
/*
	// This section checks that prom metrics are updated appropriately.
	c := make(chan prometheus.Metric, 10)

	metrics.FileCount.Collect(c)
	checkCounter(t, c, 1)

	metrics.TaskCount.Collect(c)
	checkCounter(t, c, 1)

	metrics.TestCount.Collect(c)
	checkCounter(t, c, 1)*/
}
