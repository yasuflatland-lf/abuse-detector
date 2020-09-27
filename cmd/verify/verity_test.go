package verify

import (
	"testing"
)

func TestParse(t *testing.T) {

	url := "https://www.liferay.co.jp/"
	links := []string{""}
	has, err := Parse(url, &links)

	if has == false || err == nil {
		t.Errorf("has %t error %x", has, err)
	}
}