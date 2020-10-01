package verify

import (
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
func (v *UrlScanVerifyStrategy) Request(url string) (bool, error) {
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

func (v *UrlScanVerifyStrategy) Exec(links *[]string) (bool, string, error) {
	// Check Links
	for _, link := range *links {
		ret, err := v.Request(link)

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

// Do Verification
func (v *UrlScanVerifyStrategy) Do(url string) (bool, string, error) {

	log.Info("Verification Start for <" + url + ">")

	// Parse site
	var links []string
	has, err := Scrape(url, &links)

	if has == false || err != nil {
		log.Error("Parse Error : has %t error %x", has, err)
		return has, "", err
	}

	// Check Links
	ret, link, err := v.Exec(&links)

	if err != nil {
		log.Error(err)
		return ret, link, err
	}
	//for _, link := range links {
	//	ret, err := Request(link)
	//
	//	if err != nil {
	//		log.Error(err)
	//		return false, "", err
	//	}
	//
	//	if true == ret {
	//		log.Error("Phishing link found. => %s", link)
	//		return true, link, nil
	//	} else {
	//		log.Info("OK <" + link + ">")
	//	}
	//}

	log.Info("No malicious links found.")
	return ret, link, nil
}
