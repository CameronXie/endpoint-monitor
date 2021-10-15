package monitor

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMonitor_Check(t *testing.T) {
	a := assert.New(t)
	svc := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a.Equal(http.MethodGet, r.Method)

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("resp_payload"))
	}))
	defer svc.Close()

	m := []struct {
		url      string
		hasError bool
	}{
		{
			url: svc.URL,
		},
		{
			url:      "https://127.0.0.1:1234",
			hasError: true,
		},
	}

	for _, i := range m {
		resp, err := New().Check(i.url)

		if i.hasError {
			a.NotNil(err)
			continue
		}

		zero := float64(0)
		a.Equal(i.url, resp.URL)
		a.Equal(zero, resp.DNSLookup.Seconds())
		a.Greater(resp.TCPConnTime.Seconds(), zero)
		a.GreaterOrEqual(resp.TLSHandshake.Seconds(), zero)
		a.Greater(resp.ServerTime.Seconds(), zero)
		a.Greater(resp.TotalTime.Seconds(), zero)
		a.Equal(http.StatusOK, resp.StatusCode)
		a.Nil(err)
	}
}
