package storage

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/googleapis/google-cloud-go-testing/storage/stiface"
	"github.com/m-lab/etl/etl"
	"github.com/m-lab/etl/factory"
	"github.com/m-lab/etl/row"
)

// ObjectWriter creates a writer to a named object.
// It may overwrite an existing object.
// Caller must Close() the writer, or cancel the context.
func ObjectWriter(ctx context.Context, client stiface.Client, bucket string, path string) stiface.Writer {
	b := client.Bucket(bucket)
	o := b.Object(path)
	w := o.NewWriter(ctx)
	// Set smaller chunk size to conserve memory.
	w.SetChunkSize(4 * 1024 * 1024)
	return w
}

// RowWriter implements row.Sink to a GCS file backend.
type RowWriter struct {
	ctx context.Context
	w   stiface.Writer

	bucket string
	path   string

	// These act as tokens to serialize access to the writer.
	// This allows concurrent encoding and writing, while ensuring
	// that single client access is correctly ordered.
	encoding chan struct{} // Token required for metric updates.
	writing  chan struct{} // Token required for metric updates.
}

// NewRowWriter creates a RowWriter.
func NewRowWriter(ctx context.Context, client stiface.Client, bucket string, path string) (row.Sink, error) {
	w := ObjectWriter(ctx, client, bucket, path)
	encoding := make(chan struct{}, 1)
	encoding <- struct{}{}
	writing := make(chan struct{}, 1)
	writing <- struct{}{}

	return &RowWriter{ctx: ctx, w: w, bucket: bucket, path: path, encoding: encoding, writing: writing}, nil
}

// Acquire the encoding token.
// TODO can we allow two encoders, and still sequence the writing?
func (rw *RowWriter) acquireEncodingToken() {
	<-rw.encoding
}

func (rw *RowWriter) releaseEncodingToken() {
	if len(rw.encoding) > 0 {
		log.Println("token error")
		return
	}
	rw.encoding <- struct{}{}
}

// Swap the encoding token for the write token.
// MUST already hold the write token.
func (rw *RowWriter) swapForWritingToken() {
	<-rw.writing
	rw.releaseEncodingToken()
}

func (rw *RowWriter) releaseWritingToken() {
	rw.writing <- struct{}{} // return the token.
}

// Commit commits rows, in order, to the GCS object.
// The GCS object is not available until Close is called, at which
// point the entire object becomes available atomically.
// The returned int is the number of rows written (and pending), or,
// if error is not nil, an estimate of the number of rows written.
func (rw *RowWriter) Commit(rows []interface{}, label string) (int, error) {
	rw.acquireEncodingToken()
	// First, do the encoding.  Other calls to Commit will block here
	// until encoding is done.
	// NOTE: This can cause a fairly hefty memory footprint for
	// large numbers of large rows.
	buf := bytes.NewBuffer(nil)

	for i := range rows {
		j, err := json.Marshal(rows[i])
		if err != nil {
			rw.releaseEncodingToken()
			return 0, err
		}
		buf.Write(j)
		buf.WriteByte('\n')
	}
	rw.swapForWritingToken()
	defer rw.releaseWritingToken()
	_, err := buf.WriteTo(rw.w) // This is buffered (by 4MB chunks).
	if err != nil {
		log.Println(err, rw.bucket, rw.path)
		// The caller should likely abandon the archive at this point,
		// as further writing will likely result in a corrupted file.
		return 0, err
	}

	return len(rows), nil
}

// Close synchronizes on the tokens, and closes the backing file.
func (rw *RowWriter) Close() error {
	// Take BOTH tokens, to ensure no other goroutines are still running.
	<-rw.encoding
	<-rw.writing

	close(rw.encoding)
	close(rw.writing)

	log.Println("Closing", rw.bucket, rw.path)
	err := rw.w.Close()
	if err != nil {
		log.Println(err)
	} else {
		log.Println(rw.w.Attrs())
	}
	return err
}

// SinkFactory implements factory.SinkFactory.
type SinkFactory struct {
	client       stiface.Client
	outputBucket string
}

// Get mplements factory.SinkFactory
func (sf *SinkFactory) Get(ctx context.Context, path etl.DataPath) (row.Sink, etl.ProcessingError) {
	//gcsPath := fmt.Sprintf("%s/%s/%s/%s", path.DataType, path.ExpDir, path.DatePath, path.)
	s, err := NewRowWriter(ctx, sf.client, sf.outputBucket, path.PathAndFilename()+".json")
	if err != nil {
		return nil, factory.NewError(path.DataType, "SinkFactory",
			http.StatusInternalServerError, err)
	}
	return s, nil
}

// NewSinkFactory returns the default SinkFactory
func NewSinkFactory(client stiface.Client, outputBucket string) factory.SinkFactory {
	return &SinkFactory{client: client, outputBucket: outputBucket}
}
