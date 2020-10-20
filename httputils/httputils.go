package httputils

import (
	"crypto/tls"
	"net/http"
)

// ProxyClient return a http client with "http(s)?_proxy" support
func ProxyClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			Proxy:           http.ProxyFromEnvironment,
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
}
