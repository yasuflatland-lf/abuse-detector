package verify

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestGsafeRequest(t *testing.T) {
	LoadEnv()

	cases := []struct {
		url    string
		result bool
	}{
		{url: "https://qiita.com/", result: false},
		{url: "https://loginscurrentlyattwebpage.weebly.com/", result: true},
		{url: "https://my3-uk-confirm.info", result: true},
		{url: "https://github.com/", result: false},
		{url: "https://actionukee.com/WuofvBw", result: true},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	v := NewTransparencyReportVerifyStrategy()

	for _, c := range cases {
		ret := v.Request(ctx, c.url)
		if ret.Malicious != c.result {
			t.Errorf("ret: %v result: %s}\n", ret, c.url)
		}
	}
}

func TestGsafeDo(t *testing.T) {
	LoadEnv()

	cases := []struct {
		url    string
		result bool
	}{
		{url: "https://www.google.com/", result: false},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	v := NewTransparencyReportVerifyStrategy()

	for _, c := range cases {
		ret, err := v.Do(ctx, c.url)
		if ret.Malicious != c.result || err != nil {
			t.Errorf("ret: %v result: %s}\n", ret, c.url)
		}
		fmt.Println(ret)
	}
}
