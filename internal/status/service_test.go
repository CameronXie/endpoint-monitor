package status

import (
	"context"
	"errors"
	"net"
	"testing"
	"time"

	"github.com/CameronXie/endpoint-monitor/internal/storage"
	"github.com/CameronXie/endpoint-monitor/pkg/monitor"
	"github.com/CameronXie/endpoint-monitor/pkg/scheduler"
	"github.com/stretchr/testify/assert"
)

func TestService_Monitor(t *testing.T) {
	a := assert.New(t)
	requestTime := time.Now()
	ipAddr, _ := net.ResolveIPAddr("ip", "10.0.0.1")
	c := []struct {
		ep      Endpoint
		resp    *monitor.Response
		respErr error
	}{
		{
			ep: Endpoint{
				URL:      "https://a.com",
				Interval: 10,
			},
			resp: &monitor.Response{
				URL:          "https://a.com",
				DNSLookup:    time.Nanosecond,
				TCPConnTime:  time.Millisecond,
				TLSHandshake: time.Second,
				ServerTime:   time.Minute,
				TotalTime:    time.Hour,
				RemoteAddr:   ipAddr,
				StatusCode:   200,
			},
		},
		{
			ep: Endpoint{
				URL:      "https://a.com",
				Interval: 10,
			},
			respErr: errors.New("something went wrong"),
		},
	}

	for _, i := range c {
		s := &schedulerMock{requestTime: requestTime}
		m := &monitorMock{resp: i.resp, err: i.respErr}
		store := new(storageMock)
		l := new(loggerMock)
		ctx := context.Background()
		(New(s, m, store, l)).Monitor(ctx, []Endpoint{i.ep})

		a.Equal(ctx, s.ctx)
		a.Equal(float64(i.ep.Interval), s.interval.Seconds())
		a.Equal(i.ep.URL, m.url)

		if i.respErr != nil {
			a.Equal(i.respErr.Error(), l.err)
			a.Equal(toError(i.ep.URL, i.respErr, requestTime), store.err)
			continue
		}

		a.Equal(i.resp.String(), l.info)
		a.Equal(toStatus(i.ep.URL, i.resp, requestTime), store.status)
	}
}

func TestService_Stop(t *testing.T) {
	a := assert.New(t)
	s := new(schedulerMock)
	store := new(storageMock)
	(New(s, new(monitorMock), store, new(loggerMock))).Stop()

	a.True(s.wait)
	a.True(store.flush)
}

type schedulerMock struct {
	scheduler.Scheduler
	ctx         context.Context
	interval    time.Duration
	requestTime time.Time
	wait        bool
}

func (m *schedulerMock) Schedule(ctx context.Context, d time.Duration, f scheduler.Task) {
	m.ctx = ctx
	m.interval = d
	f(m.requestTime)
}

func (m *schedulerMock) Wait() {
	m.wait = true
}

type monitorMock struct {
	monitor.Monitor
	url  string
	resp *monitor.Response
	err  error
}

func (m *monitorMock) Check(url string) (*monitor.Response, error) {
	m.url = url
	return m.resp, m.err
}

type storageMock struct {
	storage.Storage
	status *storage.Status
	err    *storage.Error
	flush  bool
}

func (m *storageMock) StoreError(err *storage.Error) {
	m.err = err
}

func (m *storageMock) StoreStatus(i *storage.Status) {
	m.status = i
}

func (m *storageMock) Flush() {
	m.flush = true
}

type loggerMock struct {
	info string
	err  string
}

func (m *loggerMock) Infoln(args ...interface{}) {
	m.info = args[0].(string)
}

func (m *loggerMock) Errorln(args ...interface{}) {
	m.err = args[0].(string)
}
