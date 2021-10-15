package monitor

import (
	"encoding/json"
	"fmt"
	"net"
	"reflect"
	"time"
)

type Response struct {
	URL          string
	DNSLookup    time.Duration
	TCPConnTime  time.Duration
	TLSHandshake time.Duration
	ServerTime   time.Duration
	TotalTime    time.Duration
	RemoteAddr   net.Addr
	StatusCode   int
}

// String returns a string representation of the Response.
func (r *Response) String() string {
	v := reflect.ValueOf(r).Elem()
	m := make(map[string]string, v.NumField())

	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		switch f.Type().String() {
		case "string", "int", "time.Duration", "net.Addr":
			if f.Interface() == nil {
				m[v.Type().Field(i).Name] = ""
				continue
			}

			m[v.Type().Field(i).Name] = fmt.Sprintf("%v", f.Interface())
		}
	}

	b, _ := json.Marshal(m)
	return string(b)
}
