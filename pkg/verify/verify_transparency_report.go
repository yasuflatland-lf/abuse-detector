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
func (v *TransparencyReportVerifyStrategy) Request(ctx context.Context, verifyUrl string) Response {
	apiUrl := os.Getenv("GOOGLE_TRANSPARENCYREPORT_API_URL")

	response := &Response{
		Result:     false,
		StatusCode: http.StatusOK,
		Error:      nil,
		Malicious:  false,
	}

	// request the HTML page.
	res, err := Fetch(apiUrl + "status?site=" + verifyUrl)

	if err != nil {
		response.StatusCode = res.StatusCode
		response.Error = err
		log.Error(err)
		return *response
	}

	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		msg := fmt.Sprintf("status code error: %d %s", res.StatusCode, res.Status)
		log.Error(msg)
		response.StatusCode = res.StatusCode
		response.Error = errors.New(msg)
		return *response
	}

	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Error(err)
		response.StatusCode = res.StatusCode
		response.Error = err
		return *response
	}

	response.Result = true
	response.StatusCode = http.StatusOK
	response.Error = nil
	response.Malicious = v.IsMalcious(string(bodyBytes))
	return *response
}

// TODO : need to make this func concurrent.
func (v *TransparencyReportVerifyStrategy) Exec(ctx context.Context, links *[]string) (bool, string, error) {

	errCh := make(chan error, len(*links))
	retCh := make(chan Result, len(*links))

	// Check Links
	for _, l := range *links {
		go func(link string) {
			retResult := &Result{}
			ret := v.Request(ctx, link)

			retResult.StatusCode = ret.StatusCode
			retResult.Error = ret.Error
			retResult.Malicious = ret.Malicious
			retResult.MaliciousLinks = append(retResult.MaliciousLinks, link)

			errCh <- ret.Error
			retCh <- *retResult
		}(l)
	}

	for _, loopTmp := range *links {
		select {
		case err := <-errCh:
			if err != nil {
				log.Error(err)
				return false, "", err
			}
		case retResult := <-retCh:
			if true == retResult.Malicious {
				log.Error("Phishing link found. => %s", retResult.MaliciousLinks[0])
				return retResult.Malicious, retResult.MaliciousLinks[0], nil
			} else {
				log.Info("OK <" + retResult.MaliciousLinks[0] + ">")
			}
		// Timeout or Cancel comes here.
		case <-ctx.Done():
			<-errCh
			return false, loopTmp, ctx.Err()
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

	// Check URL itself if it's malicious
	initRet := v.Request(ctx, url)
	result.MaliciousLinks = append(result.MaliciousLinks, url)
	result.Malicious = initRet.Malicious
	result.StatusCode = initRet.StatusCode
	if initRet.Error != nil || true == result.Malicious {
		log.Error(initRet.Error)
		return *result, initRet.Error
	}

	// Parse site
	var links []string = []string{}
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
