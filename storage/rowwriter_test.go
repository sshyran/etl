package storage_test

import (
	"context"
	"io/ioutil"
	"testing"
	"time"

	"github.com/go-test/deep"
	"github.com/googleapis/google-cloud-go-testing/storage/stiface"

	fgs "github.com/fsouza/fake-gcs-server/fakestorage"

	"github.com/m-lab/etl/storage"
)

func TestRowWriter(t *testing.T) {
	type foobar struct {
		Time time.Time
		Foo  string
		Bar  string
	}

	server := fgs.NewServer([]fgs.Object{})
	defer server.Stop()

	bucket := "fake-bucket"
	server.CreateBucket(bucket)
	c := server.Client()

	file := "foobar-file"
	rw, err := storage.NewRowWriter(context.Background(), stiface.AdaptClient(c), bucket, file)
	if err != nil {
		t.Fatal(err)
	}
	z, _ := time.LoadLocation("America/New_York")
	t1 := time.Date(1999, 1, 2, 3, 4, 5, 123456789, z)
	t2 := t1.UTC()
	rows := []interface{}{foobar{t2, "foo", "bar"}, foobar{t1, "x", "y"}}
	rw.Commit(rows, "fake-label")
	rw.Close()

	expect :=
		`{"Time":"1999-01-02T08:04:05.123456789Z","Foo":"foo","Bar":"bar"}
{"Time":"1999-01-02T03:04:05.123456789-05:00","Foo":"x","Bar":"y"}
`
	o := c.Bucket(bucket).Object(file)
	reader, err := o.NewReader(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		t.Fatal(err)
	}
	diff := deep.Equal(([]byte)(expect), data)
	if diff != nil {
		t.Error(diff)
	}
}
