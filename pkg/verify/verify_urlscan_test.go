package verify

import (
	"testing"
)

func TestRequest(t *testing.T) {
	LoadEnv()

	cases := []struct {
		url    string
		result bool
	}{
		{url: "https://vodafone-billsupport.com/", result: false},
		{url: "http://paypal-support.my-sumaya.com", result: true},
		{url: "https://my3-uk-confirm.info", result: true},
	}
	for _, c := range cases {
		ret, _ := Request(c.url)
		if ret != c.result {
			t.Errorf("ret<%t> result<%s>}\n", ret, c.url)
		}
	}
}
