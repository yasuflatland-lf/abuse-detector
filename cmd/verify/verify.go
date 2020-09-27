package verify

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
)

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

	// Find the review items
	doc.Find("").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		band := s.Find("a").Text()
		*links = append(*links, band)

		fmt.Printf("Review %d: %s \n", i, band)
	})

	if len(*links) <= 0 {
		return false, nil
	}

	return false, nil
}
