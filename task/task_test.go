// TODO(dev) add test overview
//
package task_test

import (
	"archive/tar"
	"bytes"
	"fmt"
	"reflect"
	"testing"

	"cloud.google.com/go/bigquery"

	"github.com/m-lab/etl/bq"
	"github.com/m-lab/etl/parser"
	"github.com/m-lab/etl/storage" // TODO - would be better not to have this.
	"github.com/m-lab/etl/task"
)

// Just test call to NullParser.Parse
func TestPlumbing(t *testing.T) {
	foo := [10]byte{1, 2, 3, 4, 5, 1, 2, 3, 4, 5}
	p := parser.NullParser{}
	err := p.Parse(nil, "foo", foo[:])
	if err != nil {
		fmt.Println(err)
	}
}

type NullCloser struct{}

func (nc NullCloser) Close() error {
	return nil
}

// Create a TarReader with simple test contents.
// TODO - could we break the dependency on storage here?
func MakeTestSource(t *testing.T) *storage.ETLSource {
	b := new(bytes.Buffer)
	tw := tar.NewWriter(b)
	hdr := tar.Header{Name: "foo", Mode: 0666, Typeflag: tar.TypeReg, Size: int64(8)}
	tw.WriteHeader(&hdr)
	_, err := tw.Write([]byte("biscuits"))
	if err != nil {
		t.Fatal(err)
	}

	hdr = tar.Header{Name: "bar", Mode: 0666, Typeflag: tar.TypeReg, Size: int64(11)}
	tw.WriteHeader(&hdr)
	_, err = tw.Write([]byte("butter milk"))
	if err = tw.Close(); err != nil {
		t.Fatal(err)
	}

	return &storage.ETLSource{tar.NewReader(b), NullCloser{}}
}

type TestParser struct {
	parser.Parser
	files []string
}

	// TODO - pass testName through to BQ inserter?
func (tp *TestParser) Parse(meta map[string]bigquery.Value, testName string, test []byte) error {
	tp.files = append(tp.files, testName)
	return nil
}

func (tp *TestParser) TableName() string {
	return "test"
}

// TODO(dev) - add unit tests for tgz and tar.gz files
// TODO(dev) - add good comments
func TestTarFileInput(t *testing.T) {
	rdr := MakeTestSource(t)

	var prsr TestParser
	in := bq.NullInserter{}
	tt := task.NewTask("filename", rdr, &prsr, &in, "test_table")
	fn, bb, err := tt.NextTest()
	if err != nil {
		t.Error(err)
	}
	if fn != "foo" {
		t.Error("Expected foo")
	}
	if string(bb) != "biscuits" {
		t.Error("Expected biscuits but got ", string(bb))
	}

	fn, bb, err = tt.NextTest()
	if err != nil {
		t.Error(err)
	}
	if fn != "bar" {
		t.Error("Expected bar")
	}
	if string(bb) != "butter milk" {
		t.Error("Expected butter milk but got ", string(bb))
	}

	// Reset the tar reader and create new task, to test the ProcessAllTests behavior.
	rdr = MakeTestSource(t)

	tt = task.NewTask("filename", rdr, &prsr, &in, "test_table")
	tt.ProcessAllTests()

	if len(prsr.files) != 2 {
		t.Error("Too few files ", len(prsr.files))
	}
	if !reflect.DeepEqual(prsr.files, []string{"foo", "bar"}) {
		t.Error("Not expected files: ", prsr.files)
	}

}
