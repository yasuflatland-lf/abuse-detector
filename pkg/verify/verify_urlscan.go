package verify

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
)

// For Strategy Pattern
type UrlScanVerifyStrategy struct{}

func NewUrlScanVerifyStrategy() *UrlScanVerifyStrategy {
	return &UrlScanVerifyStrategy{}
}

type UrlScanResult struct {
	Result string `json:"result"`
}

// Initial Request response
type UrlScanSubmitResponse struct {
	Results []UrlScanResult `json:"results"`
}

// UrlScanResult details
type UrlScanOverall struct {
	Malicious bool `json:"malicious"`
}

type UrlScanVerdicts struct {
	Overall UrlScanOverall `json:"overall"`
}

type UrlScanTask struct {
	URL string `json:"url"`
}

type UrlScanResultDetails struct {
	Task UrlScanTask `json:"task"`
	Verdicts UrlScanVerdicts `json:"verdicts"`
}

// Phishing site URL validation
// true if it's phishing site or false
func (v *UrlScanVerifyStrategy) IsPhishingURL(r UrlScanResult) (UrlScanResultDetails, error) {
	client := resty.New()
	emptyRet := &UrlScanResultDetails{}

	resp, err := client.R().
		EnableTrace().
		Get(r.Result)

	if err != nil {
		log.Error("Fail to read response")
		return *emptyRet, errors.Wrap(err, "Fail to read urlscan.io POST result")
	}

	var result UrlScanResultDetails
	err = json.Unmarshal([]byte(resp.String()), &result)
	if err != nil {
		log.Error(err)
		//log.Error("doc %+v", pretty.Formatter(err))
		return *emptyRet, err
	}

	//result.Verdicts.Overall.Malicious
	//result.Task.URL
	return result, nil
}

// Phishing site URL validation
// true if it's phishing site or false
func (v *UrlScanVerifyStrategy) Results(ctx context.Context, results []UrlScanResult) (bool, string, error) {
	errCh := make(chan error, len(results))
	retCh := make(chan Result, len(results))

	for _, r := range results {
		go func(result UrlScanResult) {
			retResult := &Result{}
			ret, err := v.IsPhishingURL(result)

			retResult.StatusCode = http.StatusOK
			retResult.Error = err
			retResult.Malicious = ret.Verdicts.Overall.Malicious
			retResult.MaliciousLinks = append(retResult.MaliciousLinks, ret.Task.URL)

			errCh <- err
			retCh <- *retResult
		}(r)
	}

	for _, r := range results {
		select {
		case err := <-errCh:
			if err != nil {
				log.Error(err)
				return false, "", errors.Wrap(err, "Fail to read urlscan.io POST result")
			}
		case retResult := <-retCh:
			if true == retResult.Malicious {
				log.Errorf("Phishing link found in this link. => %s", retResult.MaliciousLinks[0])
				return retResult.Malicious, retResult.MaliciousLinks[0], nil
			} else {
				if 0 != len(retResult.MaliciousLinks[0]) {
					log.Info("OK <" + retResult.MaliciousLinks[0] + ">")
				}
			}
		// Timeout or Cancel comes here.
		case <-ctx.Done():
			<-errCh
			return false, "Cancel or Done" + r.Result , ctx.Err()
		}
	}

	// Not a phishing site
	return false, "", nil
}

// Phishing site URL validation
// true if it's phishing site or false
func (v *UrlScanVerifyStrategy) Request(ctx context.Context, url string) Response {
	apiUrl := os.Getenv("URLSCAN_API_URL")
	response := &Response{
		Result:     false,
		StatusCode: http.StatusOK,
		Error:      nil,
		Malicious:  false,
	}

	hn, err := ExtractHostName(url)

	if err != nil {
		response.Error = err
		log.Error(err)
		return *response
	}

	// Create a Resty Client
	client := resty.New()

	resp, err := client.R().
		EnableTrace().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetHeader("API-Key", os.Getenv("URLSCAN_API_KEY")).
		SetQueryString("q=" + hn.HostName).
		Get(apiUrl + "/v1/search")

	if err != nil {
		log.Error("Fail to read response")
		response.StatusCode = resp.StatusCode()
		response.Error = errors.Wrap(err, "Fail to read urlscan.io POST result")
		return *response
	}

	var subRes UrlScanSubmitResponse

	err = json.Unmarshal([]byte(resp.String()), &subRes)
	if err != nil {
		log.Error(err)
		response.Error = errors.Wrap(err, "Unmarshal JSON")
		return *response
	}

	ret, _, err := v.Results(ctx, subRes.Results)

	if err != nil {
		log.Error(err)
		response.Error = err
		return *response
	}

	response.Result = true
	response.StatusCode = http.StatusOK
	response.Error = err
	response.Malicious = ret
	return *response
}

func (v *UrlScanVerifyStrategy) Exec(ctx context.Context, links *[]string) (bool, string, error) {
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
func (v *UrlScanVerifyStrategy) Do(ctx context.Context, url string) (Result, error) {

	log.Info("Verification Start for <" + url + "> - UrlScanVerifyStrategy")
	result := &Result{
		StrategyName:   "UrlScanVerifyStrategy",
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
	var links []string
	has, err := Scrape(ctx, url, &links)

	if has == false || err != nil {
		log.Error("Parse Error : has %t error %x", has, err)
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

	log.Info("No malicious links found.")
	return *result, nil
}
