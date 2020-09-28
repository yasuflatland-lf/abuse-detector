package verify

import (
	"testing"
)

func TestScrape(t *testing.T) {
	Scrape(`https://colleenpeckpdf.studio.site/`)
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
