package client

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOptionWithToken(t *testing.T) {
	token := "test-token"
	c := New(WithToken(token))
	assert.Equal(t, token, c.token)
}

func TestOptionWithHeaders(t *testing.T) {
	headers := map[string]string{
		"header-a": "value-a",
		"header-b": "value-b",
	}
	c := New(WithHeaders(headers))
	assert.Equal(t, headers, c.headers)
}

func TestOptionWithHTTPClient(t *testing.T) {
	httpclient := &http.Client{
		Timeout: time.Second * 10,
	}
	c := New(WithHTTPClient(httpclient))
	assert.Equal(t, httpclient, c.client)
}
