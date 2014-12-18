package gomock

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStub(t *testing.T) {
	transport := NewTransport()
	transport.Stub("/", func(req *http.Request) *http.Response {
		r := &http.Response{
			StatusCode: http.StatusBadRequest,
		}
		return r
	})

	if v := len(transport.layers); v != 1 {
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
	transport.Stub(".*", HandleFunc(http.StatusOK, "OK"))

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
	if resp.Body == nil {
		t.Errorf("failed")
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Failed ", err)
	}

	if string(bytes) != "OK" {
		t.Errorf("got %v body, want %v body", string(bytes), "OK")
	}
}
