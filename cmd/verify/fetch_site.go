package verify

import (
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Validate schema
// return true if the schema is https or false
func IsHttps(urlStr string) (bool, error) {
	parsedUrl, err := url.Parse(urlStr)

	if nil != err {
		log.Fatal(err)
		return false, err
	}

	return strings.EqualFold(parsedUrl.Scheme,"https"), nil
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
			MaxIdleConns:       10,
			IdleConnTimeout:    30 * time.Second,
			DisableCompression: true,
		}
		client := &http.Client{Transport: tr}
		return client.Get(url)
	} else {
		// HTTP
		return http.Get(url)
	}
}

