package verify

import (
	"context"
	"testing"
	"time"
)

func TestRequest(t *testing.T) {
	LoadEnv()

	cases := []struct {
		url    string
		result bool
	}{
		{url: "https://www.google.com/", result: false},
		{url: "http://paypal-support.my-sumaya.com", result: true},
		{url: "https://my3-uk-confirm.info", result: true},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	v := NewUrlScanVerifyStrategy()

	for _, c := range cases {
		ret, _ := v.Request(ctx, c.url)
		if ret != c.result {
			t.Errorf("ret<%t> result<%s>}\n", ret, c.url)
		}
	}
}
