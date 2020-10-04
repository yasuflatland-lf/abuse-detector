package verify

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	MaxIdleConns        = 200
	MaxIdleConnsPerHost = 200
	MaxConnsPerHost     = 200
	IdleConnTimeout     = 60 * time.Second
	DisableCompression  = true
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
// TODO : need to replace this to below.
// https://future-architect.github.io/articles/20190713/
//import (
//    "https://godoc.org/golang.org/x/net/context/ctxhttp"
//)
//
//func accessSHS(ctx context.Context) {
//    // ctxを第一引数で渡す
//    res, err := ctxhttp.Get(ctx, nil, "https://shs.sh")
//}
func Fetch(url string) (resp *http.Response, err error) {
	ret, err := IsHttps(url)

	if err != nil {
		return &http.Response{}, err
	}

	if true == ret {
		// HTTPS
		tr := &http.Transport{
			MaxIdleConns:        MaxIdleConns,
			MaxIdleConnsPerHost: MaxIdleConnsPerHost,
			MaxConnsPerHost:     MaxConnsPerHost,
			IdleConnTimeout:     IdleConnTimeout,
			DisableCompression:  DisableCompression,
			TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{Transport: tr}
		return client.Get(url)
	} else {
		// HTTP
		return http.Get(url)
	}
}
