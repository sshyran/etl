package schema

import (
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/m-lab/annotation-service/api"
	"github.com/m-lab/go/bqx"
	"github.com/m-lab/ndt-server/ndt5/c2s"
	"github.com/m-lab/ndt-server/ndt5/control"
	"github.com/m-lab/ndt-server/ndt5/s2c"
)

type NDTResult struct {
	// GitShortCommit is the Git commit (short form) of the running server code.
	GitShortCommit string
	// Version is the symbolic version (if any) of the running server code.
	Version string

	// All data members should all be self-describing. In the event of confusion,
	// rename them to add clarity rather than adding a comment.
	ServerIP   string
	ServerPort int
	ClientIP   string
	ClientPort int

	StartTime time.Time
	EndTime   time.Time

	// ndt5
	Control *control.ArchivalData `json:",omitempty"`
	C2S     *c2s.ArchivalData     `json:",omitempty"`
	S2C     *s2c.ArchivalData     `json:",omitempty"`

	// ndt7
	// Upload   *model.ArchivalData `json:",omitempty"`
	// Download *model.ArchivalData `json:",omitempty"`
}

// NDTResultRow defines the BQ schema for the data.NDTResult produced by the
// ndt-server for NDT client measurements.
type NDTResultRow struct {
	ParseInfo *ParseInfo
	TestID    string    `json:"test_id,string" bigquery:"test_id"`
	LogTime   int64     `json:"log_time,int64" bigquery:"log_time"`
	Result    NDTResult `json:"result" bigquery:"result"`
}

// Schema returns the BigQuery schema for NDTResultRow.
func (row *NDTResultRow) Schema() (bigquery.Schema, error) {
	sch, err := bigquery.InferSchema(row)
	if err != nil {
		return bigquery.Schema{}, err
	}
	docs := FindSchemaDocsFor(row)
	for _, doc := range docs {
		bqx.UpdateSchemaDescription(sch, doc)
	}
	rr := bqx.RemoveRequired(sch)
	return rr, err
}

// Implement row.Annotatable
// This is a trivial implementation, as the schema does not yet include
// annotations, and probably will not until we integrate UUID Annotator.
func (row *NDTResultRow) GetLogTime() time.Time {
	return time.Now()
}
func (row *NDTResultRow) GetClientIPs() []string {
	return []string{}
}
func (row *NDTResultRow) GetServerIP() string {
	return ""
}
func (row *NDTResultRow) AnnotateClients(map[string]*api.Annotations) error {
	return nil
}
func (row *NDTResultRow) AnnotateServer(*api.Annotations) error {
	return nil
}
