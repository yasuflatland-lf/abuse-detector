package verify

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/url"
)

// Check if the URL includes schema
// true if it does or false
func IsSchema(urlStr string) (bool, error) {
	parsedUrl, err := url.Parse(urlStr)

	if nil != err {
		log.Fatal(err)
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
	// Request the HTML page.
	res, err := Fetch(url)
	if err != nil {
		log.Fatal(err)
		return false, err
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		msg := fmt.Sprintf("status code error: %d %s", res.StatusCode, res.Status)
		log.Fatal(msg)
		return false, errors.New(msg)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
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

	return true, nil
}
