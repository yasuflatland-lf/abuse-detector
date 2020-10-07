package verify

import (
	"context"
	"crypto/tls"
	"golang.org/x/net/context/ctxhttp"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
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
func Fetch(ctx context.Context,  url string) (resp *http.Response, err error) {
	MaxIdleConns, _ := strconv.Atoi(os.Getenv("COMMON_MAX_IDLE_CONNS"))
	MaxIdleConnsPerHost, _ := strconv.Atoi(os.Getenv("COMMON_MAX_IDLE_CONN_SPER_HOST"))
	MaxConnsPerHost, _ := strconv.Atoi(os.Getenv("COMMON_MAX_CONNS_PER_HOST"))
	IdleConnTimeout, _ := strconv.Atoi(os.Getenv("COMMON_IDLE_CONN_TIMEOUT"))
	DisableCompression, _ := strconv.ParseBool(os.Getenv("COMMON_DISABLE_COMPRESSION"))
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
			IdleConnTimeout:     time.Duration(IdleConnTimeout) * time.Second,
			DisableCompression:  DisableCompression,
			TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{Transport: tr}
		return ctxhttp.Get(ctx, client, url)
	} else {
		// HTTP
		return ctxhttp.Get(ctx, http.DefaultClient, url)
	}
}
