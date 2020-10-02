package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"studio.design/studio-abuse-detector/pkg/verify"

	"github.com/stretchr/testify/assert"
)

func TestHelloHandler(t *testing.T) {
	verify.LoadEnv()

	router := NewRouter()

	req := httptest.NewRequest("GET", "/verify?url=https://bono760lbk.site/", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "{\"strategyName\":\"\",\"link\":[],\"malicious\":false,\"error\":null}\n", rec.Body.String())
}
