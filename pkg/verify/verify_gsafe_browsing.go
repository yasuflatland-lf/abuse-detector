package verify

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
)

type GsafeClient struct {
	ClientID      string `json:"clientId"`
	ClientVersion string `json:"clientVersion"`
}

type GsafeThreatEntries struct {
	URL string `json:"url"`
}

type GsafeThreatInfo struct {
	ThreatTypes      []string             `json:"threatTypes"`
	PlatformTypes    []string             `json:"platformTypes"`
	ThreatEntryTypes []string             `json:"threatEntryTypes"`
	ThreatEntries    []GsafeThreatEntries `json:"threatEntries"`
}

type GsafeRequestBody struct {
	Client     GsafeClient     `json:"client"`
	ThreatInfo GsafeThreatInfo `json:"threatInfo"`
}

type GsafeThreat struct {
	URL string `json:"url"`
}

type GsafeThreatEntryMetadata struct {
}

type GsafeMatches struct {
	ThreatType          string                   `json:"threatType"`
	PlatformType        string                   `json:"platformType"`
	ThreatEntryType     string                   `json:"threatEntryType"`
	Threat              GsafeThreat              `json:"threat"`
	ThreatEntryMetadata GsafeThreatEntryMetadata `json:"threatEntryMetadata"`
	CacheDuration       string                   `json:"cacheDuration"`
}

type GsafeResponseBody struct {
	Matches []GsafeMatches `json:"matches"`
}

func ConvertToGsafeThreatEntries(urls []string, entries *[]GsafeThreatEntries) {
	for _, url := range urls {
		*entries = append(*entries, GsafeThreatEntries{URL: url})
	}
	return
}

// Phishing site URL validation
// true if it's phishing site or false
func GsafeRequest(urls []string) (bool, error) {
	apiUrl := os.Getenv("GOOGLE_SAFE_BROWSING_API_URL")
	apiKey := os.Getenv("GOOGLE_SAFE_BROWSING_API_KEY")

	var entries = []GsafeThreatEntries{}
	ConvertToGsafeThreatEntries(urls, &entries)

	// Create a Resty Client
	client := resty.New()

	resp, err := client.R().
		EnableTrace().
		SetHeader("Content-Type", "application/json").
		SetBody(GsafeRequestBody{
			Client: GsafeClient{
				ClientID:      "Studio Inc",
				ClientVersion: "1.5.2",
			},
			ThreatInfo: GsafeThreatInfo{
				ThreatTypes: []string{
					"THREAT_TYPE_UNSPECIFIED",
					"MALWARE",
					"SOCIAL_ENGINEERING",
					"UNWANTED_SOFTWARE",
					"POTENTIALLY_HARMFUL_APPLICATION",
				},
				PlatformTypes:    []string{"ANY_PLATFORM", "PLATFORM_TYPE_UNSPECIFIED"},
				ThreatEntryTypes: []string{"URL", "EXECUTABLE", "THREAT_ENTRY_TYPE_UNSPECIFIED"},
				ThreatEntries:    entries,
			},
		}).
		SetQueryString("key=" + apiKey).
		Post(apiUrl)

	if err != nil {
		log.Error("Fail to read response")
		return false, errors.Wrap(err, "Fail to read Google Safe Browsing API POST result")
	}

	var subRes GsafeResponseBody

	err = json.Unmarshal([]byte(resp.String()), &subRes)
	if err != nil {
		log.Error(err)
		return false, nil
	}
	fmt.Print(subRes)

	return true, nil
}
