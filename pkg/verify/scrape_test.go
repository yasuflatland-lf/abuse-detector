package verify

import (
	"testing"
)

func TestScrape(t *testing.T) {
	var links = []string{}
	url := `https://www.google.com/`
	ret, err := Scrape(url, &links)

	if err != nil || true != ret {
		t.Errorf("url <%s>, %v", url, err)
	}

	if len(links) <= 0 {
		t.Error("The length of links is invalid.")
	}
}

func TestFindHref(t *testing.T) {
	cases := []struct {
		hrefs  []string
		result bool
		retStr string
	}{
		{
			hrefs:  []string{"https://vodafone-billsupport.com/"},
			result: true,
			retStr: "https://vodafone-billsupport.com/",
		},
		{
			hrefs:  []string{"href"},
			result: false,
			retStr: "",
		},
		{
			hrefs:  []string{"href", "https://vodafone-billsupport.com/"},
			result: true,
			retStr: "https://vodafone-billsupport.com/",
		},
		{
			hrefs:  []string{"https://vodafone-billsupport.com/", "href", "http://example.com"},
			result: true,
			retStr: "https://vodafone-billsupport.com/",
		},
	}
	for _, c := range cases {
		ret, stat := FindHref(c.hrefs)
		if stat != c.result || ret != c.retStr {
			t.Errorf("ret<%s> stat<%t> should fetch <%s>\n", ret, c.result, c.retStr)
		}
	}
}
