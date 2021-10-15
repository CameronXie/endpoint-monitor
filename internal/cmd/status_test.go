package cmd

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/CameronXie/endpoint-monitor/internal/status"
	"github.com/stretchr/testify/assert"
)

func Test_statusRun(t *testing.T) {
	a := assert.New(t)
	pwd, _ := os.Getwd()
	f := filepath.Join(pwd, "..", "..", "tests", "cmd", "config.yml")
	s := make(chan os.Signal, 1)
	svc := new(svcMock)
	l := new(loggerMock)

	go func() {
		s <- os.Interrupt
	}()
	statusRun(&f, s, svc, l)(nil, nil)

	a.Equal([]status.Endpoint{{
		URL:      "http://localhost",
		Interval: 10,
	}}, svc.eps)
	a.True(svc.stop)
	a.Equal(stopMonitorMsg, l.info)
}

type svcMock struct {
	eps  []status.Endpoint
	stop bool
}

func (m *svcMock) Monitor(_ context.Context, eps []status.Endpoint) {
	m.eps = eps
}

func (m *svcMock) Stop() {
	m.stop = true
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
