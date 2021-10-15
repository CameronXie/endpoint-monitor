package status

import (
	"context"
	"time"

	"github.com/CameronXie/endpoint-monitor/internal/storage"
	"github.com/CameronXie/endpoint-monitor/pkg/monitor"
	"github.com/CameronXie/endpoint-monitor/pkg/scheduler"
)

type MonitorService interface {
	Monitor(ctx context.Context, eps []Endpoint)
	Stop()
}

type Logger interface {
	Infoln(args ...interface{})
	Errorln(args ...interface{})
}

type Endpoint struct {
	URL      string
	Interval uint
}

type service struct {
	scheduler scheduler.Scheduler
	monitor   monitor.Monitor
	storage   storage.Storage
	logger    Logger
}

func (s *service) Monitor(ctx context.Context, eps []Endpoint) {
	for _, i := range eps {
		e := i
		s.scheduler.Schedule(ctx, time.Duration(e.Interval)*time.Second, func(t time.Time) {
			resp, err := s.monitor.Check(e.URL)
			if err != nil {
				s.logger.Errorln(err.Error())
				s.storage.StoreError(toError(e.URL, err, t))
				return
			}

			s.logger.Infoln(resp.String())
			s.storage.StoreStatus(toStatus(e.URL, resp, t))
		})
	}
}

func (s *service) Stop() {
	s.scheduler.Wait()
	s.storage.Flush()
}

func toError(url string, err error, reqTime time.Time) *storage.Error {
	return &storage.Error{
		Endpoint:    url,
		Message:     err.Error(),
		RequestTime: reqTime,
	}
}

func toStatus(url string, resp *monitor.Response, reqTime time.Time) *storage.Status {
	return &storage.Status{
		Endpoint:     url,
		DNSLookup:    resp.DNSLookup.Nanoseconds(),
		TCPConnTime:  resp.TCPConnTime.Nanoseconds(),
		TLSHandshake: resp.TLSHandshake.Nanoseconds(),
		ServerTime:   resp.ServerTime.Nanoseconds(),
		TotalTime:    resp.TotalTime.Nanoseconds(),
		RemoteAddr:   resp.RemoteAddr.String(),
		StatusCode:   resp.StatusCode,
		RequestTime:  reqTime,
	}
}

func New(
	s scheduler.Scheduler,
	m monitor.Monitor,
	store storage.Storage,
	l Logger,
) MonitorService {
	return &service{
		scheduler: s,
		monitor:   m,
		storage:   store,
		logger:    l,
	}
}
