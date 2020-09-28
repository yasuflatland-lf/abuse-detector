package verify

import (
	"encoding/json"
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
)

// Phishing site URL validation
// true if it's phishing site or false
func IsPhishingURL(r Result) (bool, error) {
	client := resty.New()

	resp, err := client.R().
		EnableTrace().
		Get(r.Result)

	if err != nil {
		log.Error("Fail to read response")
		return false, errors.Wrap(err, "Fail to read urlscan.io POST result")
	}

	var result ResultDetails
	err = json.Unmarshal([]byte(resp.String()), &result)
	if err != nil {
		log.Error(err)
		//log.Error("doc %+v", pretty.Formatter(err))
		return false, nil
	}

	//fmt.Println("RESULT :", result.Verdicts.Overall.Malicious)
	return result.Verdicts.Overall.Malicious, nil
}

// Phishing site URL validation
// true if it's phishing site or false
func Results(results []Result) (bool, error) {
	for _, r := range results {
		ret, err := IsPhishingURL(r)

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
func Request(url string) (bool, error) {
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

	var subRes submitResponse

	err = json.Unmarshal([]byte(resp.String()), &subRes)
	if err != nil {
		log.Error(err)
		//log.Error("doc %+v", pretty.Formatter(err))
		return false, nil
	}

	return Results(subRes.Results)
}

// Run Verification
func Run(url string) (bool, string, error) {

	log.Info("Verification Start for <" + url + ">")

	// Parse site
	links := []string{url}
	has, err := Parse(url, &links)

	if has == false || err != nil {
		log.Error("Parse Error : has %t error %x", has, err)
		return has, "", err
	}

	// Check Links
	for _, link := range links {
		ret, err := Request(link)

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

	log.Error("No malicious links found.")
	return false, "", nil
}

type Result struct {
	Result string `json:"result"`
}

// Initial Request response
type submitResponse struct {
	Results []Result `json:"results"`
}

// Result details
type Overall struct {
	Malicious bool `json:"malicious"`
}

type Verdicts struct {
	Overall Overall `json:"overall"`
}

type ResultDetails struct {
	Verdicts Verdicts `json:"verdicts"`
}
