package gomock

import (
	"net"
	"net/http"
	"time"
)

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
