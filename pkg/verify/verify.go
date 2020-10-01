package verify

import (
	"net/url"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("verify")

var logFmt = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} PID=%{pid} MOD=%{module} PKG=%{shortpkg} %{shortfile} FUNC=%{shortfunc} ▶ %{level:.4s} %{id:03x} %{color:reset} %{message}`,
)

// Verify Interface
type Verify interface {
	Request(url string) (bool, error)
	Do(url string) (bool, string, error)
}

type VerifyExecutor struct {
	url      string
	strategy Verify
}

func NewVerifyExecutor(url string, v Verify) *VerifyExecutor {
	return &VerifyExecutor{
		url:      url,
		strategy: v,
	}
}

type HostNames struct {
	URL      string
	HostName string
}

// Extract valid URL for verification API
// Return URL with either http or https or return empty string
func ExtractHostName(urlStr string) (HostNames, error) {
	hn := &HostNames{
		URL:      "",
		HostName: "",
	}

	u, err := url.Parse(urlStr)

	if err != nil {
		log.Error(err)
		return *hn, err
	}

	isSchema, err := IsSchema(urlStr)

	if err != nil {
		log.Error(err)
		return *hn, err
	}

	if u.Hostname() != "" && true == isSchema {
		hn.URL = u.Scheme + "://" + u.Hostname()
		hn.HostName = u.Hostname()
	}

	return *hn, nil
}

// Check if the URL includes schema
// true if it does or false
func IsSchema(urlStr string) (bool, error) {
	parsedUrl, err := url.Parse(urlStr)

	if nil != err {
		log.Error(err)
		return false, err
	}

	var bSchema bool = true
	if len(parsedUrl.Scheme) == 0 {
		// No schema
		bSchema = false
	} else if parsedUrl.Scheme != "http" && parsedUrl.Scheme != "https" {
		// Neither http nor https
		bSchema = false
	}

	return bSchema, nil
}
