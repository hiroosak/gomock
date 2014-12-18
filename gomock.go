// gomock package is mocking http request server.
package gomock

import (
	"fmt"
	"net"
	"net/http"
	"regexp"
	"time"
)

type Handle func(req *http.Request) *http.Response

type Layer struct {
	Pattern    *regexp.Regexp
	HandleFunc Handle
}

type Transport struct {
	Transport http.RoundTripper
	layers    []Layer
}

func NewTransport() *Transport {
	return &Transport{
		Transport: http.DefaultTransport,
		layers:    []Layer{},
	}
}

func (t *Transport) RegisterProtocol(scheme string, rt http.RoundTripper) {
}

func (t *Transport) CloseIdleConnections() {
}

func (t *Transport) CancelRequest(req *http.Request) {
}

// Stub returns handle function.
func (t *Transport) Stub(m interface{}, handle Handle) error {
	l, err := newLayer(m, handle)
	if err != nil {
		return err
	}
	t.layers = append(t.layers, l)
	return nil
}

func newLayer(m interface{}, handle Handle) (Layer, error) {
	switch v := m.(type) {
	case *regexp.Regexp:
		return Layer{
			Pattern:    v,
			HandleFunc: handle,
		}, nil
	case string:
		return Layer{
			Pattern:    regexp.MustCompile(v),
			HandleFunc: handle,
		}, nil
	}
	return Layer{}, fmt.Errorf("invalid m %v", m)
}

func (t *Transport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	for _, l := range t.layers {
		if l.Pattern.MatchString(req.URL.String()) {
			resp := l.HandleFunc(req)
			if resp != nil {
				return resp, err
			}
		}
	}
	return t.Transport.RoundTrip(req)
}

// HandleFunc returns a mock response object.
func HandleFunc(statusCode int, body string) Handle {
	return func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: statusCode,
			Body:       NewReadCloser(body),
			Request:    req,
		}
	}
}

// DefaultTransport set to DefaultRansport.
func SetDefaultTransport(t *Transport) {
	http.DefaultTransport = t
}

// ResetDefaultTransport clear http.DefaultTransport.
func ResetDefaultTransport() {
	http.DefaultTransport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 10 * time.Second,
	}
}
