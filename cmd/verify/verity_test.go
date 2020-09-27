package verify

import (
	"github.com/kr/pretty"
	"testing"
)

func TestParse(t *testing.T) {

	url := "https://www.liferay.co.jp/"
	links := []string{""}
	has, err := Parse(url, &links)

	if has == false || err != nil {
		t.Errorf("has %t error %x", has, err)
	}

	t.Logf("links %+v", pretty.Formatter(links))
}