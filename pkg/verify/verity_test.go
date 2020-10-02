package verify

import (
	"testing"
)

func TestExtractHostName(t *testing.T) {
	cases := []struct {
		url    string
		result string
	}{
		{url: "https://www.liferay.co.jp/?q=aaa", result: "https://www.liferay.co.jp"},
		{url: "http://violet-evergarden.jp/aaa", result: "http://violet-evergarden.jp"},
		{url: "/some/path", result: ""},
		{url: "smb://some/path", result: ""},
	}

	for _, c := range cases {
		ret, err := ExtractHostName(c.url)

		if err != nil || ret.URL != c.result {
			t.Errorf("Url %s is error. should be %s", ret, c.result)
		}
	}
}

func TestIsSchema(t *testing.T) {
	cases := []struct {
		url    string
		result bool
	}{
		{url: "https://www.liferay.co.jp/", result: true},
		{url: "http://violet-evergarden.jp/", result: true},
		{url: "/some/path", result: false},
		{url: "smb://some/path", result: false},
	}

	for _, c := range cases {
		ret, err := IsSchema(c.url)

		if err != nil || ret != c.result {
			t.Errorf("Url %s is error. should be %t", c.url, c.result)
		}
	}
}
