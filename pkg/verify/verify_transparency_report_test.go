package verify

import (
	"testing"
)

func TestGsafeRequest(t *testing.T) {
	LoadEnv()

	cases := []struct {
		url    string
		result bool
	}{
		{url: "https://qiita.com/", result: false},
		{url: "https://vodafone-billsupport.com/", result: true},
		{url: "http://paypal-support.my-sumaya.com", result: true},
		{url: "https://my3-uk-confirm.info", result: true},
		{url: "https://github.com/", result: false},
		{url: "https://actionukee.com/WuofvBw", result: true},
	}

	v := NewTransparencyReportVerifyStrategy()

	for _, c := range cases {
		ret, _ := v.Request(c.url)
		if ret != c.result {
			t.Errorf("ret: %t result: %s}\n", ret, c.url)
		}
	}
}
