// This files contains schema for Paris TraceRoute tests.
package schema

import (
	"context"
	"log"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/m-lab/annotation-service/api"
	"github.com/m-lab/etl/metrics"
	"github.com/m-lab/go/bqx"
	"github.com/m-lab/go/rtx"
	"github.com/prometheus/client_golang/prometheus"
)

// TODO(dev): use mixed case Go variable names throughout

type ParisTracerouteHop struct {
	Protocol         string            `json:"protocal,string"`
	Src_ip           string            `json:"src_ip,string"`
	Src_af           int32             `json:"src_af,int32"`
	Dest_ip          string            `json:"dest_ip,string"`
	Dest_af          int32             `json:"dest_af,int32"`
	Src_hostname     string            `json:"src_hostname,string"`
	Dest_hostname    string            `json:"dest_hostname,string"`
	Rtt              []float64         `json:"rtt,[]float64"`
	Src_geolocation  api.GeolocationIP `json:"src_geolocation"`
	Dest_geolocation api.GeolocationIP `json:"dest_geolocation"`
}

type MLabConnectionSpecification struct {
	Server_ip          string            `json:"server_ip,string"`
	Server_af          int32             `json:"server_af,int32"`
	Client_ip          string            `json:"client_ip,string"`
	Client_af          int32             `json:"client_af,int32"`
	Data_direction     int32             `json:"data_direction,int32"`
	Server_geolocation api.GeolocationIP `json:"server_geolocation"`
	Client_geolocation api.GeolocationIP `json:"client_geolocation"`
}

// PT describes a single BQ row of PT data.
type PT struct {
	TestID               string                      `json:"test_id,string" bigquery:"test_id"`
	Project              int32                       `json:"project,int32" bigquery:"project"`
	TaskFilename         string                      `json:"task_filename,string" bigquery:"task_filename"`
	ParseTime            time.Time                   `json:"parse_time" bigquery:"parse_time"`
	ParserVersion        string                      `json:"parser_version,string" bigquery:"parser_version"`
	LogTime              int64                       `json:"log_time,int64" bigquery:"log_time"`
	Connection_spec      MLabConnectionSpecification `json:"connection_spec"`
	Paris_traceroute_hop ParisTracerouteHop          `json:"paris_traceroute_hop"`
	Type                 int32                       `json:"type,int32"`
}

func assertTCPRowIsValueSaver(r *PTRow) {
	func(bigquery.ValueSaver) {}(r)
}

func init() {
	var err error
	ptSchema, err = (&PT{}).Schema()
	rtx.Must(err, "Error generating PT schema")
}

var ptSchema bigquery.Schema

// Save implements bigquery.ValueSaver
func (row *PT) Save() (map[string]bigquery.Value, string, error) {
	ss := bigquery.StructSaver{Schema: ptSchema, Struct: row}
	return ss.Save()
}

// Schema returns the Bigquery schema for PT.
func (row *PT) Schema() (bigquery.Schema, error) {
	sch, err := bigquery.InferSchema(row)
	if err != nil {
		return bigquery.Schema{}, err
	}
	rr := bqx.RemoveRequired(sch)
	return rr, nil
}

// Implement parser.Annotatable

// GetLogTime returns the timestamp that should be used for annotation.
func (row *PT) GetLogTime() time.Time {
	return row.LogTime
}

// GetClientIPs returns the client (remote) IP for annotation.  See parser.Annotatable
func (row *PT) GetClientIPs() []string {
	clientIPs := make([]string, 1)
	clientIPs = append(clientIPs, row.Connection_spec.Client_ip)
	return clientIPs
}

// GetServerIP returns the server (local) IP for annotation.  See parser.Annotatable
func (row *TCPRow) GetServerIP() string {
	return row.Connection_spec.Server_ip
}

func (row *PT) AnnotateClients(annMap map[string]*api.Annotations) error {
	ip := row.Connection_spec.Client_ip
	ann, ok := annMap[ip]
	if !ok {
		metrics.AnnotationMissingCount.WithLabelValues("PT: No annotation for client IP").Inc()
		return nil
	}
	if ann.Geo == nil {
		metrics.AnnotationMissingCount.WithLabelValues("PT: Empty client ann.Geo").Inc()
	} else {
		row.Connection_spec.Client_geolocation = ann.Geo
	}
  // TODO: add ASN to PT schema
	return nil
}

func (row *PT) AnnotateServer(local *api.Annotations) error {
	if local == nil {
		return nil
	}
	row.Connection_spec.Server_geolocation. = local.Geo
  // TODO: add ASN to PT schema
	return nil
}

func (row *PT) AnnotateHops(annMap map[string]*api.Annotations) error {
	for _, hop := range row.Paris_traceroute_hop {
		annSrc, ok := annMap[hop.Src_ip]
		if !ok {
			metrics.AnnotationMissingCount.WithLabelValues("PT: No annotation for hop src IP").Inc()
		}
		if annSrc.Geo == nil {
			metrics.AnnotationMissingCount.WithLabelValues("PT: Empty hop src ann.Geo").Inc()
		} else {
			hop.Src_geolocation = annSrc.Geo
		}

		annDest, ok := annMap[hop.Dest_ip]
		if !ok {
			metrics.AnnotationMissingCount.WithLabelValues("PT: No annotation for hop dest IP").Inc()
		}
		if annDest.Geo == nil {
			metrics.AnnotationMissingCount.WithLabelValues("PT: Empty hop dest ann.Geo").Inc()
		} else {
			hop.Dest_geolocation = annDest.Geo
		}
	}
  return nil
}

// AnnotatePT sent one batch requests for all hops in one PT test.
func (row *PT) AnnotatePT(requestIPs []string, logTime time.Time) error {
	response, err := buf.ann.GetAnnotations(context.Background(), logTime, requestIPs, "PT")
	if err != nil {
		log.Println("error in PT GetAnnotations: ", err)
		metrics.AnnotationErrorCount.With(prometheus.
			Labels{"source": "PT: RPC err in GetAnnotations."}).Inc()
		return err
	}
	annMap := response.Annotations
	if annMap == nil {
		log.Println("empty client annotation response")
		metrics.AnnotationErrorCount.With(prometheus.
			Labels{"source": "Client IP: empty response"}).Inc()
		return ErrAnnotationError
	}

	row.AnnotateClients(annMap)

	ann, ok := annMap[row.Connection_spec.Server_ip]
	if !ok {
		metrics.AnnotationMissingCount.WithLabelValues("PT: No annotation for server IP").Inc()
		return nil
	}
	if ann.Geo == nil {
		metrics.AnnotationMissingCount.WithLabelValues("PT: Empty server ann.Geo").Inc()
	} else {
		row.AnnotateServer(ann)
	}

	row.AnnotateHops(annMap)
	return nil
}
