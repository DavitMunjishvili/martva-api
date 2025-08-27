package httpclient

import (
	"crypto/tls"
	"net/http"
)

// New creates a new HTTP client with InsecureSkipVerify set to true.
func New() *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return &http.Client{Transport: tr}
}
