package storage

import (
	"fmt"
	"testing"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"github.com/stretchr/testify/assert"
)

func TestNewInfluxDBStorage(t *testing.T) {
	a := assert.New(t)
	org, bucket, url, token := "org_test", "bucket_test", "example.com", "token"
	m := []struct {
		batchSize uint
		expected  uint
	}{
		{
			batchSize: 0,
			expected:  defaultBatchSize,
		},
		{
			batchSize: 15,
			expected:  15,
		},
	}

	for _, i := range m {
		s := NewInfluxDBStorage(&InfluxdbStorageConfig{
			URL:       url,
			Token:     token,
			Org:       org,
			Bucket:    bucket,
			BatchSize: i.batchSize,
		}).(*influxdbStorage)

		a.Equal(org, s.org)
		a.Equal(bucket, s.bucket)
		a.Equal(url, s.client.ServerURL())
		a.Equal(fmt.Sprintf("Token %v", token), s.client.HTTPService().Authorization())
		a.Equal(i.expected, s.client.Options().BatchSize())
	}
}

func TestInfluxdbStorage_Store(t *testing.T) {
	a := assert.New(t)
	org, bucket := "org", "bucket"
	writeTime := time.Now()
	m := []struct {
		status   *Status
		expected *write.Point
	}{
		{
			status: &Status{
				Endpoint:     "example.com",
				DNSLookup:    1,
				TCPConnTime:  2,
				TLSHandshake: 3,
				ServerTime:   4,
				TotalTime:    5,
				RemoteAddr:   "10.0.0.1:443",
				StatusCode:   200,
				RequestTime:  writeTime,
			},
			expected: influxdb2.NewPoint(
				statusMeasurement,
				map[string]string{
					"Endpoint":      "example.com",
					"RemoteAddress": "10.0.0.1:443",
				},
				map[string]interface{}{
					"DNSLookup":    1,
					"TCPConnTime":  2,
					"TLSHandshake": 3,
					"ServerTime":   4,
					"TotalTime":    5,
					"StatusCode":   200,
				},
				writeTime,
			),
		},
	}

	for _, i := range m {
		apiMock := &writeAPIMock{
			writeTime: writeTime,
		}
		s := &influxdbStorage{
			org:    org,
			bucket: bucket,
			client: &influxdbClientMock{
				writeAPIMock: apiMock,
			},
		}

		s.StoreStatus(i.status)
		a.Equal(i.expected, apiMock.point)
	}
}

func TestInfluxdbStorage_StoreError(t *testing.T) {
	a := assert.New(t)
	org, bucket := "org_err", "bucket_err"
	writeTime := time.Now()
	m := []struct {
		err      *Error
		expected *write.Point
	}{
		{
			err: &Error{
				Endpoint:    "example.com",
				Message:     "connection error",
				RequestTime: writeTime,
			},
			expected: influxdb2.NewPoint(
				errorMeasurement,
				map[string]string{
					"Endpoint": "example.com",
				},
				map[string]interface{}{
					"ErrorMessage": "connection error",
				},
				writeTime,
			),
		},
	}

	for _, i := range m {
		apiMock := &writeAPIMock{
			writeTime: writeTime,
		}
		s := &influxdbStorage{
			org:    org,
			bucket: bucket,
			client: &influxdbClientMock{
				writeAPIMock: apiMock,
			},
		}

		s.StoreError(i.err)
		a.Equal(i.expected, apiMock.point)
	}
}

func TestInfluxdbStorage_Flush(t *testing.T) {
	a := assert.New(t)
	apiMock := new(writeAPIMock)
	(&influxdbStorage{
		client: &influxdbClientMock{
			writeAPIMock: apiMock,
		},
	}).Flush()

	a.True(apiMock.isFlushed)
}

type influxdbClientMock struct {
	influxdb2.Client
	writeAPIMock *writeAPIMock
}

func (m *influxdbClientMock) WriteAPI(_, _ string) api.WriteAPI {
	return m.writeAPIMock
}

type writeAPIMock struct {
	api.WriteAPI
	point     *write.Point
	writeTime time.Time
	isFlushed bool
}

func (m *writeAPIMock) WritePoint(point *write.Point) {
	m.point = point
}

func (m *writeAPIMock) Flush() {
	m.isFlushed = true
}
