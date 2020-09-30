package verify

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
)

// There's no official documentation exposed for transparencyreport.
// This definitions are based on the response from the API v3
// Could be changed without a notice as this does not look like exposed API.
const (
	errorFlag1Idx   = 1
	errorFlag1Value = "2"
	errorFlag2Idx   = 4
	errorFlag2Value = "1"
)

func GsafeGetResponse(respStr string) []string {
	// Clean up response
	var noNLstr string = strings.ReplaceAll(string(respStr), "\n", "")
	r := regexp.MustCompile(`\[\[(\S+)\]\]`)
	result := r.FindAllStringSubmatch(noNLstr, -1)

	return strings.Split(result[0][1], ",")
}

func GsafeIsMalcious(respStr string) bool {
	resp := GsafeGetResponse(respStr)
	if resp[errorFlag1Idx] == errorFlag1Value &&
		resp[errorFlag2Idx] == errorFlag2Value {
		return true
	}
	return false
}

// Referred https://transparencyreport.google.com/safe-browsing/search
func GsafeRequest(verifyUrl string) (bool, error) {
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

	return GsafeIsMalcious(string(bodyBytes)), nil
}

func Do(verifyUrl string) (bool, error) {
	return true, nil
}
