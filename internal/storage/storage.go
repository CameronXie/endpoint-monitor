package storage

import (
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

const (
	defaultBatchSize  = 50
	statusMeasurement = "status"
	errorMeasurement  = "error"
)

type Status struct {
	Endpoint     string
	DNSLookup    int64
	TCPConnTime  int64
	TLSHandshake int64
	ServerTime   int64
	TotalTime    int64
	RemoteAddr   string
	StatusCode   int
	RequestTime  time.Time
}

type Error struct {
	Endpoint    string
	Message     string
	RequestTime time.Time
}

type Storage interface {
	// StoreStatus stores *Status.
	StoreStatus(status *Status)

	// StoreError stores *Error
	StoreError(err *Error)

	// Flush forces all pending items to be stored.
	Flush()
}

type InfluxdbStorageConfig struct {
	URL       string
	Token     string
	Org       string
	Bucket    string
	BatchSize uint
}

type influxdbStorage struct {
	org    string
	bucket string
	client influxdb2.Client
}

func (s *influxdbStorage) StoreStatus(status *Status) {
	s.client.WriteAPI(s.org, s.bucket).WritePoint(
		influxdb2.NewPoint(
			statusMeasurement,
			map[string]string{
				"Endpoint":      status.Endpoint,
				"RemoteAddress": status.RemoteAddr,
			},
			map[string]interface{}{
				"DNSLookup":    status.DNSLookup,
				"TCPConnTime":  status.TCPConnTime,
				"TLSHandshake": status.TLSHandshake,
				"ServerTime":   status.ServerTime,
				"TotalTime":    status.TotalTime,
				"StatusCode":   status.StatusCode,
			},
			status.RequestTime,
		),
	)
}

func (s *influxdbStorage) StoreError(err *Error) {
	s.client.WriteAPI(s.org, s.bucket).WritePoint(
		influxdb2.NewPoint(
			errorMeasurement,
			map[string]string{
				"Endpoint": err.Endpoint,
			},
			map[string]interface{}{
				"ErrorMessage": err.Message,
			},
			err.RequestTime,
		),
	)
}

func (s *influxdbStorage) Flush() {
	s.client.WriteAPI(s.org, s.bucket).Flush()
}

func NewInfluxDBStorage(config *InfluxdbStorageConfig) Storage {
	if config.BatchSize == 0 {
		config.BatchSize = defaultBatchSize
	}

	return &influxdbStorage{
		org:    config.Org,
		bucket: config.Bucket,
		client: influxdb2.NewClientWithOptions(
			config.URL,
			config.Token,
			influxdb2.DefaultOptions().SetBatchSize(config.BatchSize),
		),
	}
}
