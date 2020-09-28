package verify

import (
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	MaxIdleConns       = 10
	IdleConnTimeout    = 30 * time.Second
	DisableCompression = true
)

// Validate schema
// return true if the schema is https or false
func IsHttps(urlStr string) (bool, error) {
	parsedUrl, err := url.Parse(urlStr)

	if nil != err {
		log.Error(err)
		return false, err
	}

	return strings.EqualFold(parsedUrl.Scheme, "https"), nil
}

// Fetch URL response
// Automatically detect https or http
func Fetch(url string) (resp *http.Response, err error) {
	ret, err := IsHttps(url)

	if err != nil {
		return &http.Response{}, err
	}

	if true == ret {
		// HTTPS
		tr := &http.Transport{
			MaxIdleConns:       MaxIdleConns,
			IdleConnTimeout:    IdleConnTimeout,
			DisableCompression: DisableCompression,
		}
		client := &http.Client{Transport: tr}
		return client.Get(url)
	} else {
		// HTTP
		return http.Get(url)
	}
}
