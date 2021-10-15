package monitor

import (
	"github.com/go-resty/resty/v2"
)

type Monitor interface {
	// Check the given url / endpoints, returns *Response or error if there is any.
	Check(url string) (*Response, error)
}

type monitor struct {
	client *resty.Client
}

func (m *monitor) Check(url string) (*Response, error) {
	resp, err := m.client.R().EnableTrace().Get(url)

	if err != nil {
		return nil, err
	}

	ti := resp.Request.TraceInfo()
	return &Response{
		URL:          url,
		DNSLookup:    ti.DNSLookup,
		TCPConnTime:  ti.TCPConnTime,
		TLSHandshake: ti.TLSHandshake,
		ServerTime:   ti.ServerTime,
		TotalTime:    ti.TotalTime,
		RemoteAddr:   ti.RemoteAddr,
		StatusCode:   resp.StatusCode(),
	}, nil
}

func New() Monitor {
	return &monitor{client: resty.New()}
}
