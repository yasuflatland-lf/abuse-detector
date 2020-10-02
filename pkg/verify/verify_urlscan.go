package verify

import (
	"context"
	"encoding/json"
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

type UrlScanResultDetails struct {
	Verdicts UrlScanVerdicts `json:"verdicts"`
}

// Phishing site URL validation
// true if it's phishing site or false
func (v *UrlScanVerifyStrategy) IsPhishingURL(r UrlScanResult) (bool, error) {
	client := resty.New()

	resp, err := client.R().
		EnableTrace().
		Get(r.Result)

	if err != nil {
		log.Error("Fail to read response")
		return false, errors.Wrap(err, "Fail to read urlscan.io POST result")
	}

	var result UrlScanResultDetails
	err = json.Unmarshal([]byte(resp.String()), &result)
	if err != nil {
		log.Error(err)
		//log.Error("doc %+v", pretty.Formatter(err))
		return false, err
	}

	return result.Verdicts.Overall.Malicious, nil
}

// Phishing site URL validation
// true if it's phishing site or false
func (v *UrlScanVerifyStrategy) Results(results []UrlScanResult) (bool, error) {
	for _, r := range results {
		ret, err := v.IsPhishingURL(r)

		if err != nil {
			log.Error("Fail to read response")
			return false, errors.Wrap(err, "Fail to read urlscan.io POST result")
		}

		if true == ret {
			// Phishing site detected. Return right away
			return true, nil
		}
	}

	// Not a phishing site
	return false, nil
}

// Phishing site URL validation
// true if it's phishing site or false
func (v *UrlScanVerifyStrategy) Request(ctx context.Context, url string) (bool, error) {
	apiUrl := os.Getenv("URLSCAN_API_URL")

	hn, err := ExtractHostName(url)

	if err != nil {
		log.Error(err)
		return false, err
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
		return false, errors.Wrap(err, "Fail to read urlscan.io POST result")
	}

	var subRes UrlScanSubmitResponse

	err = json.Unmarshal([]byte(resp.String()), &subRes)
	if err != nil {
		log.Error(err)
		//log.Error("doc %+v", pretty.Formatter(err))
		return false, nil
	}

	return v.Results(subRes.Results)
}

func (v *UrlScanVerifyStrategy) Exec(ctx context.Context, links *[]string) (bool, string, error) {
	errCh := make(chan error, len(*links))
	retCh := make(chan Result, len(*links))

	// Check Links
	for _, l := range *links {
		go func(link string) {
			retResult := &Result{}
			ret, err := v.Request(ctx, link)

			retResult.Malicious = ret
			retResult.MaliciousLinks = append(retResult.MaliciousLinks, link)

			errCh <- err
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

	log.Info("Verification Start for <" + url + ">")
	result := &Result{
		StrategyName:   "UrlScanVerifyStrategy",
		Malicious:      false,
		MaliciousLinks: []string{},
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
