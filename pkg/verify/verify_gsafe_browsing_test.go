package verify

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertToGsafeThreatEntries(t *testing.T) {
	entries := []GsafeThreatEntries{}
	urls := []string{"a", "b", "c", "d"}
	ConvertToGsafeThreatEntries(urls, &entries)

	assert.Equal(t, len(urls), len(entries))
}

func TestGsafeRequest(t *testing.T) {
	LoadEnv()

	urls := []string{
		"https://www.facebook.com/",
		"https://yourordercheckout.com/diepost/ch",
		"http://www.verificar-bcpmovil.cobra-ks.com",
		"https://bigattmm.weebly.com/",
	}
	ret, err := GsafeRequest(urls)

	if err != nil {
		log.Error("Fail to read response")
		t.Errorf("ret<%t> result<%v>}\n", ret, urls)
	}
}
