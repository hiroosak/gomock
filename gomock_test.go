package gomock

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStub(t *testing.T) {
	Stub("/", func(req *http.Request) *http.Response {
		r := &http.Response{
			Status:     "200 OK",
			StatusCode: http.StatusOK,
		}
		return r
	})

	if v := len(DefaultTransport.layers); v != 1 {
		t.Errorf("len(DefaultTransport) == 1, but %v", v)
	}
}

func TestRoundTrip(t *testing.T) {
	handler := func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	transport := NewTransport()
	transport.Stub("/", func(req *http.Request) *http.Response {
		r := &http.Response{
			Status:     "200 OK",
			StatusCode: http.StatusOK,
		}
		return r
	})

	client := &http.Client{
		Transport: transport,
	}

	resp, err := client.Get(server.URL)
	if err != nil {
		t.Errorf("Failed ", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("got %v status code, want %v status code", resp.StatusCode, http.StatusOK)
	}
}
