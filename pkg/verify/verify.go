package verify

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
	"github.com/op/go-logging"
	"github.com/thoas/go-funk"
)

var log = logging.MustGetLogger("verify")

var logFmt = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} PID=%{pid} MOD=%{module} PKG=%{shortpkg} %{shortfile} FUNC=%{shortfunc} ▶ %{level:.4s} %{id:03x} %{color:reset} %{message}`,
)

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

// Parse HTML
func Parse(url string, links *[]string) (bool, error) {

	hn, err := ExtractHostName(url)

	if err != nil {
		log.Error(err)
		return false, err
	}

	// request the HTML page.
	res, err := Fetch(hn.URL)
	if err != nil {
		log.Error(err)
		return false, err
	}

	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		msg := fmt.Sprintf("status code error: %d %s", res.StatusCode, res.Status)
		log.Error(msg)
		return false, errors.New(msg)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Error(err)
		return false, err
	}

	// TODO Need to extract this as function parameter to be able to pass multiple link extract methods
	// Find the review items
	duplicateCheck := make(map[string]struct{})
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		attr, exists := s.Attr("href")

		// Check links only for schemas, http or https are included
		isSchema, _ := IsSchema(attr)

		if _, ok := duplicateCheck[attr]; !ok &&
			true == exists &&
			true == isSchema {

			duplicateCheck[attr] = struct{}{}
			*links = append(*links, attr)
		}
	})

	// Unique strings
	*links = funk.UniqString(*links)

	return true, nil
}
