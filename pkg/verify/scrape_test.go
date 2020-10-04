package verify

import (
	"context"
	"testing"
	"time"
)

func TestScrape(t *testing.T) {
	LoadEnv()

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	var links = []string{}
	url := `https://www.google.com/`
	ret, err := Scrape(ctx, url, &links)

	if err != nil || true != ret {
		t.Errorf("Url <%s>, %v", url, err)
	}

	if len(links) <= 0 {
		t.Error("The length of links is invalid.")
	}
}

func TestFindHref(t *testing.T) {
	LoadEnv()

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
