package monitor

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestResponse_String(t *testing.T) {
	a := assert.New(t)
	ipAddr, _ := net.ResolveIPAddr("ip", "10.0.0.1")
	m := []struct {
		resp     *Response
		expected string
	}{
		{
			resp: &Response{
				URL:          "https://localhost",
				DNSLookup:    time.Nanosecond,
				TCPConnTime:  time.Millisecond,
				TLSHandshake: time.Second,
				ServerTime:   time.Minute,
				TotalTime:    time.Hour,
				RemoteAddr:   ipAddr,
				StatusCode:   200,
			},
			//nolint: lll
			expected: `{"DNSLookup":"1ns","RemoteAddr":"10.0.0.1","ServerTime":"1m0s","StatusCode":"200","TCPConnTime":"1ms","TLSHandshake":"1s","TotalTime":"1h0m0s","URL":"https://localhost"}`,
		},
		{
			resp: &Response{
				URL:          "https://localhost",
				DNSLookup:    time.Nanosecond,
				TCPConnTime:  time.Millisecond,
				TLSHandshake: time.Second,
				ServerTime:   time.Minute,
				TotalTime:    time.Hour,
				StatusCode:   200,
			},
			//nolint: lll
			expected: `{"DNSLookup":"1ns","RemoteAddr":"","ServerTime":"1m0s","StatusCode":"200","TCPConnTime":"1ms","TLSHandshake":"1s","TotalTime":"1h0m0s","URL":"https://localhost"}`,
		},
	}

	for _, i := range m {
		a.Equal(i.expected, i.resp.String())
	}
}
