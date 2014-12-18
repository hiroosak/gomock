package gomock

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"regexp"
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

func TestDefaultTransport(t *testing.T) {
	handler := func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	transport := NewTransport()
	transport.Stub(".*", HandleFunc(http.StatusOK, "OK"))

	SetDefaultTransport(transport)
	defer ResetDefaultTransport()

	client := &http.Client{}

	resp, err := client.Get(server.URL)
	if err != nil {
		t.Errorf("Failed ", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("got %v status code, want %v status code", resp.StatusCode, http.StatusOK)
	}
}

func ExampleStub() {
	transport := NewTransport()
	transport.Stub(`http://example.jp`, HandleFunc(http.StatusOK, "hello"))

	client := &http.Client{
		Transport: transport,
	}
	resp, err := client.Get(`http://example.jp`)
	if err != nil {
		log.Fatalln("error:", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln("error:", err)
	}
	fmt.Println(string(body))
	// Output: hello
}

func ExampleRegexpStub() {
	transport := NewTransport()
	transport.Stub(regexp.MustCompile(`/foo$`), HandleFunc(http.StatusOK, "foo"))
	transport.Stub(regexp.MustCompile(`/bar$`), HandleFunc(http.StatusOK, "bar"))

	client := &http.Client{
		Transport: transport,
	}
	resp, err := client.Get(`http://example.jp/foo`)
	if err != nil {
		log.Fatalln("error:", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln("error:", err)
	}
	fmt.Println(string(body))
	// Output: foo
}
