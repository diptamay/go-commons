package helpers

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"
)

type HTTPClientSettings struct {
	Connect          time.Duration
	ConnKeepAlive    time.Duration
	ExpectContinue   time.Duration
	IdleConn         time.Duration
	MaxAllIdleConns  int
	MaxHostIdleConns int
	ResponseHeader   time.Duration
	TLSHandshake     time.Duration
	TLSConfig        *tls.Config
}

func NewHTTPClient(httpSettings HTTPClientSettings) *http.Client {
	tr := &http.Transport{
		ResponseHeaderTimeout: httpSettings.ResponseHeader,
		Proxy:                 http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			KeepAlive: httpSettings.ConnKeepAlive,
			Timeout:   httpSettings.Connect,
		}).DialContext,
		MaxIdleConns:          httpSettings.MaxAllIdleConns,
		IdleConnTimeout:       httpSettings.IdleConn,
		TLSHandshakeTimeout:   httpSettings.TLSHandshake,
		MaxIdleConnsPerHost:   httpSettings.MaxHostIdleConns,
		ExpectContinueTimeout: httpSettings.ExpectContinue,
	}

	if &httpSettings.TLSConfig != nil {
		tr.TLSClientConfig = httpSettings.TLSConfig
	}

	return &http.Client{
		Transport: tr,
	}
}

func NewHTTPClientRecommended() *http.Client {
	return NewHTTPClient(NewHTTPClientSettingsWithDefaults())
}

func NewHTTPClientInsecure() *http.Client {

	settings := NewHTTPClientSettingsWithDefaults()
	settings.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	return NewHTTPClient(settings)
}

func NewHTTPClientSettingsWithDefaults() HTTPClientSettings {
	return HTTPClientSettings{
		Connect:          5 * time.Second,
		ExpectContinue:   1 * time.Second,
		IdleConn:         90 * time.Second,
		ConnKeepAlive:    30 * time.Second,
		MaxAllIdleConns:  100,
		MaxHostIdleConns: 10,
		ResponseHeader:   5 * time.Second,
		TLSHandshake:     5 * time.Second,
	}
}
