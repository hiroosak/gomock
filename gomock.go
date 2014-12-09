// gomock package is mocking http request server.
package gomock

import "net/http"

type Layer func(req *http.Request) (resp *http.Response, err error)

type Transport struct {
	Transport http.RoundTripper
	Layers    []Layer
}

func NewTransport() *Transport {
	return &Transport{
		Transport: http.DefaultTransport,
		Layers:    []Layer{},
	}
}

var DefaultTransport = NewTransport()

func (t *Transport) AddLayer(handler Layer) {
	t.Layers = append(t.Layers, handler)
}

func (t *Transport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	for _, l := range t.Layers {
		resp, err := l(req)
		if resp != nil || err != nil {
			return resp, err
		}
	}
	return t.Transport.RoundTrip(req)
}
