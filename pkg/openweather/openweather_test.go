package openweather

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var ctx = context.Background()

func setup(opts ...Opt) (*Client, *http.ServeMux, func()) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	opts = append([]Opt{WithBaseURL(server.URL)}, opts...)
	client, _ := New("key123", opts...)

	return client, mux, server.Close
}

func readFileContents(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}

	return string(bytes), err
}

func testClientServices(t *testing.T, c *Client) {
	services := []string{
		"Current",
	}

	cp := reflect.ValueOf(c)
	cv := reflect.Indirect(cp)

	for _, s := range services {
		require.Falsef(t, cv.FieldByName(s).IsNil(), "c.%s should not be nil", s)
	}
}

func TestNewClient(t *testing.T) {
	client, err := New("")
	require.NoError(t, err)
	testClientServices(t, client)
}

func TestNewClientFromEnv(t *testing.T) {
	_, err := NewFromEnv()
	require.EqualError(t, err, "OPEN_WEATHER_API_KEY environment variable is not set: cannot initialize client")

	os.Setenv("OPEN_WEATHER_API_KEY", "key123")
	defer os.Unsetenv("OPEN_WEATHER_API_KEY")

	client, err := NewFromEnv()
	require.NoError(t, err)
	require.Equal(t, "key123", client.apiKey)
	testClientServices(t, client)
}

func TestNewClient_Error(t *testing.T) {
	errorOpt := func(c *Client) error {
		return errors.New("foo")
	}

	_, err := New("", errorOpt)
	require.EqualError(t, err, "foo")
}

func TestClient_NewRequest_Errors(t *testing.T) {
	client, _, teardown := setup()
	defer teardown()

	_, err := client.NewRequest(" ", "/foo", nil)
	assert.EqualError(t, err, `net/http: invalid method " "`)

	_, err = client.NewRequest(http.MethodGet, ":", nil)
	urlErr, ok := err.(*url.Error)
	assert.True(t, ok)
	assert.Equal(t, "parse", urlErr.Op)
}

func TestClient_ErrorResponse(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v1/test", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)

		err := r.ParseForm()
		require.NoError(t, err)

		require.Equal(t, defaultUserAgent, r.Header.Get(headerUserAgent))
		require.Equal(t, "key123", r.Form.Get(paramAPIKey))

		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{
			"message": "an error occurred"
		}`)
	})

	req, err := client.NewRequest(http.MethodGet, "api/v1/test", nil)
	require.NoError(t, err)

	resp, err := client.Do(ctx, req, nil)
	require.IsType(t, &ErrorResponse{}, err)
	require.EqualError(t, err, fmt.Sprintf("GET %s/api/v1/test?appid=key123: 400 an error occurred", client.BaseURL))
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
}
