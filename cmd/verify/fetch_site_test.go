package verify

import (
	"reflect"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/kr/pretty"
)

func TestIsHttps(t *testing.T) {
	cases := []struct {
		url  string
		result bool
	}{
		{url: "https://www.liferay.co.jp/", result: true},
		{url: "http://violet-evergarden.jp/", result: false},
	}
	for _, c := range cases {
		ret, _ := IsHttps(c.url)
		if ret != c.result {
			t.Errorf("ret<%t> result<%s>}\n",ret, c.url)
		}
	}
}

func TestFetchSiteSmoke(t *testing.T) {
	cases := []struct {
		url  string
	}{
		{url: "https://www.liferay.co.jp/"},
		{url: "http://violet-evergarden.jp/"},
	}
	for _, c := range cases {
		doc, err := Fetch(c.url)

		expectDoc := goquery.Document{}

		if true == reflect.DeepEqual(&expectDoc, &doc) || nil != err {
			t.Errorf("doc %+v error %x", doc, err)
		}
		t.Logf("doc %+v", pretty.Formatter(doc))
	}
}