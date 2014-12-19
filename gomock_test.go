package gomock

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestStub(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/me", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(200)
		fmt.Fprintf(w, "OK")
	})

	transport := NewTransport()
	transport.Stub("graph.facebook.com", mux)

	client := &http.Client{
		Transport: transport,
	}

	resp, err := client.Get("https://graph.facebook.com/me")
	if err != nil {
		t.Errorf(err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("should return %v. But %v", http.StatusOK, resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf(err.Error())
	}

	if string(body) != "OK" {
		t.Errorf("should return http.StatusOK. But %v", string(body))
	}
}

func TestDefaultTransport(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/me", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(200)
		fmt.Fprintf(w, "OK")
	})

	transport := NewTransport()
	transport.Stub(".*", mux)
	SetDefaultTransport(transport)
	defer ResetDefaultTransport()

	client := &http.Client{}

	resp, err := client.Get("http://example.jp/me")
	if err != nil {
		t.Errorf("Failed ", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("got %v status code, want %v status code", resp.StatusCode, http.StatusOK)
	}
}
