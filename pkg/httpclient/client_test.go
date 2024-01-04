package httpclient

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// RoundTripFunc .
type RoundTripFunc func(req *http.Request) *http.Response

// RoundTrip .
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func TestOgameClient_Do(t *testing.T) {
	c := Client{userAgent: "test", Client: &http.Client{Transport: RoundTripFunc(func(req *http.Request) *http.Response {
		// Test request parameters
		return &http.Response{
			StatusCode: 200,
			// Send response to be tested
			Body: ioutil.NopCloser(bytes.NewBufferString(`OK`)),
			// Must be set to non-nil value or it panics
			Header: make(http.Header),
		}
	})}}
	req, _ := http.NewRequest(http.MethodGet, "http://test.com", nil)
	_, err := c.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, "test", req.Header.Get("User-Agent"))
}

func TestOgameClient_SetUserAgent(t *testing.T) {
	c := Client{userAgent: "test", Client: &http.Client{Transport: RoundTripFunc(func(req *http.Request) *http.Response {
		// Test request parameters
		return &http.Response{
			StatusCode: 200,
			// Send response to be tested
			Body: ioutil.NopCloser(bytes.NewBufferString(`OK`)),
			// Must be set to non-nil value or it panics
			Header: make(http.Header),
		}
	})}}
	req, _ := http.NewRequest(http.MethodGet, "http://test.com", nil)
	_, err := c.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, "test", req.Header.Get("User-Agent"))
}
