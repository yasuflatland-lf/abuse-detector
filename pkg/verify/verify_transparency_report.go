package verify

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
)

// For Strategy Pattern
type TransparencyReportVerifyStrategy struct{}

func NewTransparencyReportVerifyStrategy() *TransparencyReportVerifyStrategy {
	return &TransparencyReportVerifyStrategy{}
}

// There's no official documentation exposed for transparencyreport.
// This definitions are based on the response from the API v3
// Could be changed without a notice as this does not look like exposed API.
const (
	errorFlag1Idx   = 1
	errorFlag1Value = "2"
	errorFlag2Idx   = 4
	errorFlag2Value = "1"
)

func (v *TransparencyReportVerifyStrategy) Response(respStr string) []string {
	// Clean up response
	var noNLstr string = strings.ReplaceAll(string(respStr), "\n", "")
	r := regexp.MustCompile(`\[\[(\S+)\]\]`)
	result := r.FindAllStringSubmatch(noNLstr, -1)

	return strings.Split(result[0][1], ",")
}

func (v *TransparencyReportVerifyStrategy) IsMalcious(respStr string) bool {
	resp := v.Response(respStr)
	if resp[errorFlag1Idx] == errorFlag1Value &&
		resp[errorFlag2Idx] == errorFlag2Value {
		return true
	}
	return false
}

// Referred https://transparencyreport.google.com/safe-browsing/search
func (v *TransparencyReportVerifyStrategy) Request(ctx context.Context, verifyUrl string) (bool, error) {
	apiUrl := os.Getenv("GOOGLE_TRANSPARENCYREPORT_API_URL")

	// request the HTML page.
	res, err := Fetch(apiUrl + "status?site=" + verifyUrl)

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

	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Error(err)
		return false, err
	}

	return v.IsMalcious(string(bodyBytes)), nil
}

// TODO : need to make this func concurrent.
func (v *TransparencyReportVerifyStrategy) Exec(ctx context.Context, links *[]string) (bool, string, error) {
	// Check Links
	for _, link := range *links {
		ret, err := v.Request(ctx, link)

		if err != nil {
			log.Error(err)
			return false, "", err
		}

		if true == ret {
			log.Error("Phishing link found. => %s", link)
			return true, link, nil
		} else {
			log.Info("OK <" + link + ">")
		}
	}

	return false, "", nil
}

// TODO : Refactor this to common func with Template?
// Do Verification
func (v *TransparencyReportVerifyStrategy) Do(ctx context.Context, url string) (Result, error) {

	log.Info("Verification Start for <" + url + ">")
	result := &Result{
		StrategyName:   "TransparencyReportVerifyStrategy",
		Malicious:      false,
		MaliciousLinks: []string{},
	}

	// Parse site
	var links []string = []string{url}
	has, err := Scrape(ctx, url, &links)

	if has == false || err != nil {
		log.Errorf("Parse Error : result <%t>", has)
		return *result, err
	}

	// Check Links
	ret, link, err := v.Exec(ctx, &links)
	result.MaliciousLinks = append(result.MaliciousLinks, link)
	result.Malicious = ret

	if err != nil {
		log.Error(err)
		return *result, err
	}

	return *result, nil
}
