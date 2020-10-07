package openweather

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWithHTTPClient(t *testing.T) {
	_, err := New("", WithHTTPClient(nil))
	require.EqualError(t, err, "*http.Client: cannot be nil")

	httpClient := &http.Client{}
	client, err := New("", WithHTTPClient(httpClient))
	require.NoError(t, err)
	require.Same(t, httpClient, client.client)
}

func TestWithUserAgent(t *testing.T) {
	c, err := New("", WithUserAgent("test"))
	require.NoError(t, err)
	require.Equal(t, "test", c.UserAgent)
}

func TestWithBaseURL(t *testing.T) {
	c, err := New("", WithBaseURL(":"))
	urlErr, ok := err.(*url.Error)
	require.True(t, ok)
	require.Equal(t, "parse", urlErr.Op)

	baseURL := "http://localhost:8080"
	c, err = New("", WithBaseURL(baseURL))
	require.NoError(t, err)
	require.Equal(t, baseURL, c.BaseURL.String())
}
