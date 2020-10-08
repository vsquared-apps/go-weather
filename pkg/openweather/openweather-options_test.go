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
	client, err := New("", WithUserAgent("test"))
	require.NoError(t, err)
	require.Equal(t, "test", client.UserAgent)
}

func TestWithBaseURL(t *testing.T) {
	client, err := New("", WithBaseURL(":"))
	require.Error(t, err)
	urlErr, ok := err.(*url.Error)
	require.True(t, ok)
	require.Equal(t, "parse", urlErr.Op)

	baseURL := "http://localhost:8080"
	client, err = New("", WithBaseURL(baseURL))
	require.NoError(t, err)
	require.Equal(t, baseURL, client.BaseURL.String())
}

func TestUseUnits(t *testing.T) {
	client, err := New("")
	require.NoError(t, err)
	require.Empty(t, client.units)

	req, err := client.NewRequest(http.MethodGet, "api/v1/test", nil)
	require.NoError(t, err)
	require.Empty(t, req.URL.Query().Get("units"))

	client, err = New("", UseMetricUnits)
	require.NoError(t, err)
	require.Equal(t, "metric", client.units)

	req, err = client.NewRequest(http.MethodGet, "api/v1/test", nil)
	require.NoError(t, err)
	require.Equal(t, "metric", req.URL.Query().Get("units"))

	client, err = New("", UseImperialUnits)
	require.NoError(t, err)
	require.Equal(t, "imperial", client.units)

	req, err = client.NewRequest(http.MethodGet, "api/v1/test", nil)
	require.NoError(t, err)
	require.Equal(t, "imperial", req.URL.Query().Get("units"))
}
