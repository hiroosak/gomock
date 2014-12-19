// gomock package is mocking http request server.
package gomock

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"sync"
)

var mutex *sync.Mutex = &sync.Mutex{}

// Layer includes a multiplexer and domain pattern.
type Layer struct {
	Pattern *regexp.Regexp
	Mux     *http.ServeMux
}

// Transport implements http.RoundTripper.
type Transport struct {
	Transport http.RoundTripper
	layers    []Layer
}

// NewTransport returns a new Transport, with matching url.
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

// Stub returns mux function.
func (t *Transport) Stub(m interface{}, mux *http.ServeMux) error {
	l, err := newLayer(m, mux)
	if err != nil {
		return err
	}
	t.layers = append(t.layers, l)
	return nil
}

func newLayer(m interface{}, mux *http.ServeMux) (Layer, error) {
	switch v := m.(type) {
	case *regexp.Regexp:
		return Layer{
			Pattern: v,
			Mux:     mux,
		}, nil
	case string:
		return Layer{
			Pattern: regexp.MustCompile(v),
			Mux:     mux,
		}, nil
	}
	return Layer{}, fmt.Errorf("invalid m %v", m)
}

func (t *Transport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	mutex.Lock()
	defer mutex.Unlock()

	for _, l := range t.layers {
		if l.Pattern.MatchString(req.URL.String()) {
			server := httptest.NewServer(l.Mux)
			defer server.Close()

			newReq, err := url.Parse(server.URL)
			if err != nil {
				return nil, err
			}
			newReq.Path = req.URL.Path
			newReq.RawQuery = req.URL.Query().Encode()
			req.URL = newReq
		}
	}
	return t.Transport.RoundTrip(req)
}
