// gomock package is mocking http request server.
package gomock

import (
	"fmt"
	"net/http"
	"regexp"
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

func (t *Transport) Stub(m interface{}, handle Handle) error {
	l, err := newLayer(m, handle)
	if err != nil {
		return err
	}
	t.layers = append(DefaultTransport.layers, l)
	return nil
}

var DefaultTransport = NewTransport()

func Stub(m interface{}, handle Handle) error {
	l, err := newLayer(m, handle)
	if err != nil {
		return err
	}
	DefaultTransport.layers = append(DefaultTransport.layers, l)
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
		resp := l.HandleFunc(req)
		if resp != nil {
			return resp, err
		}
	}
	return t.Transport.RoundTrip(req)
}
