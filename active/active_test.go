// The active package provides code for managing processing of an entire
// directory of task files.
package active_test

import (
	"context"
	"errors"
	"log"
	"strings"
	"sync"
	"testing"
	"time"

	"golang.org/x/sync/semaphore"

	"cloud.google.com/go/storage"
	"github.com/m-lab/go/cloudtest"

	"github.com/m-lab/etl/active"
)

func init() {
	// Always prepend the filename and line number.
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

type counter struct {
	lock    sync.Mutex
	t       *testing.T
	fail    int
	success int
}

func (c *counter) processTask(tf *active.TaskFile) error {
	time.Sleep(10 * time.Millisecond)
	c.lock.Lock()
	defer c.lock.Unlock()
	if !strings.HasPrefix(tf.Path(), "gs://foobar/") {
		c.t.Error("Invalid path:", tf.Path())
	}
	if c.fail > 0 {
		log.Println("Intentional temporary failure:", tf)
		c.fail--
		return errors.New("intentional test failure")
	}
	log.Println(tf)
	c.success++
	return nil
}

type TokenSource struct {
	sem *semaphore.Weighted
}

func (ts *TokenSource) Acquire(ctx context.Context) error {
	return ts.sem.Acquire(ctx, 1)
}
func (ts *TokenSource) Release() {
	ts.sem.Release(1)
}

func TestProcessAll(t *testing.T) {
	fc := cloudtest.GCSClient{}
	fc.AddTestBucket("foobar",
		cloudtest.BucketHandle{
			ObjAttrs: []*storage.ObjectAttrs{
				&storage.ObjectAttrs{Bucket: "foobar", Name: "ndt/2019/01/01/obj1", Updated: time.Now()},
				&storage.ObjectAttrs{Bucket: "foobar", Name: "ndt/2019/01/01/obj2", Updated: time.Now()},
				&storage.ObjectAttrs{Bucket: "foobar", Name: "ndt/2019/01/01/obj3"},
				&storage.ObjectAttrs{Bucket: "foobar", Name: "ndt/2019/01/01/subdir/obj4", Updated: time.Now()},
				&storage.ObjectAttrs{Bucket: "foobar", Name: "ndt/2019/01/01/subdir/obj5", Updated: time.Now()},
				&storage.ObjectAttrs{Bucket: "foobar", Name: "obj6", Updated: time.Now()},
			}})

	// First four attempts will fail.  This means that one of the 3 tasks will have two failures.
	p := counter{t: t, fail: 4}
	// Retry once per file.  This means one of the 3 tasks will never succeed.
	fs, err := active.NewFileSource(fc, "fake", "gs://foobar/ndt/2019/01/01/")
	if err != nil {
		t.Fatal(err)
	}
	disp := active.NewDispatcher(fs, p.processTask, 1)
	tokens := TokenSource{semaphore.NewWeighted(1)}
	wg, err := disp.ProcessAll(context.Background(), fs, &tokens)
	if err != nil {
		t.Fatal(err)
	}
	wg.Wait()

	// At this point, we may be still draining the last tasks.
	//	tokens.Acquire(context.Background())
	//tokens.Release()

	// One file should have failed twice.  Others should have failed once, then succeeded.
	if len(disp.Errors()) != 1 {
		t.Errorf("ProcessAll() had %d errors %v, %v", len(disp.Errors()), disp.Errors()[0], disp.Errors())
	}

	if p.success != 2 {
		t.Error("Expected 3 successes, got", p.success)
	}
}

func TestNoFiles(t *testing.T) {
	fc := cloudtest.GCSClient{}
	fc.AddTestBucket("foobar",
		cloudtest.BucketHandle{
			ObjAttrs: []*storage.ObjectAttrs{}})

	p := counter{t: t} // All processing attempts will succeed.
	fs, err := active.NewFileSource(fc, "fake", "gs://foobar/ndt/2019/01/01/")
	if err != nil {
		t.Fatal(err)
	}
	disp := active.NewDispatcher(fs, p.processTask, 1)
	tokens := TokenSource{semaphore.NewWeighted(2)}
	wg, err := disp.ProcessAll(context.Background(), fs, &tokens)
	if err != nil {
		t.Fatal(err)
	}
	wg.Wait()
	// At this point, we may be still draining the last tasks.
	// Hacky way to hopefully drain all the tasks.
	/*	tokens.Acquire(context.Background())
		tokens.Acquire(context.Background())
		tokens.Release()
		tokens.Release()*/

	if len(disp.Errors()) > 0 {
		t.Error("ProcessAll() had errors", disp.Errors())
	}

	// processTask should never be called, because there are no files.
	if p.success+p.fail != 0 {
		t.Error("Expected 0 successes, got", p.success)
	}
}
