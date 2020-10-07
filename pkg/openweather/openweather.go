package openweather

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

const (
	defaultBaseURL   = "https://api.openweathermap.org"
	defaultUserAgent = "github.com/vsquared-apps/go-weather/pkg/openweather"

	mediaTypeJSON = "application/json"

	headerUserAgent   = "User-Agent"
	headerContentType = "Content-Type"
	headerAccept      = "Accept"

	paramAPIKey = "appid"
)

// Client makes requests to the OpenWeather API.
type Client struct {
	// HTTP client used to make requests.
	client *http.Client

	BaseURL   *url.URL
	UserAgent string

	apiKey string

	Current *CurrentService
}

func newClient(apiKey string) *Client {
	baseURL, _ := url.Parse(defaultBaseURL)

	client := &Client{
		client:    &http.Client{},
		BaseURL:   baseURL,
		UserAgent: defaultUserAgent,
		apiKey:    apiKey,
	}

	client.Current = &CurrentService{client: client}

	return client
}

// New returns a new OpenWeather client.
func New(apiKey string, opts ...Opt) (*Client, error) {
	client := newClient(apiKey)

	for _, opt := range opts {
		if err := opt(client); err != nil {
			return nil, err
		}
	}

	return client, nil
}

// NewFromEnv returns a new OpenWeather client with the API key set to the value of the
// OPEN_WEATHER_API_KEY environment variable.
func NewFromEnv(opts ...Opt) (*Client, error) {
	apiKey, ok := os.LookupEnv("OPEN_WEATHER_API_KEY")
	if !ok {
		return nil, errors.New("OPEN_WEATHER_API_KEY environment variable is not set: cannot initialize client")
	}
	return New(apiKey, opts...)
}

// NewRequest creates an API request. A relative URL can be provided in urlStr, which will be resolved to the
// BaseURL of the Client. Relative URLS should always be specified without a preceding slash. If specified, the
// value pointed to by body is JSON encoded and included in as the request body.
func (c *Client) NewRequest(method string, path string, body interface{}) (*http.Request, error) {
	u, err := c.BaseURL.Parse(path)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	if body != nil {
		err = json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	reqBody := bytes.NewReader(buf.Bytes())
	req, err := http.NewRequest(method, u.String(), reqBody)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add(paramAPIKey, c.apiKey)

	req.URL.RawQuery = q.Encode()
	req.Header.Add(headerUserAgent, c.UserAgent)
	req.Header.Add(headerContentType, mediaTypeJSON)
	req.Header.Add(headerAccept, mediaTypeJSON)

	return req, nil
}

// Do sends an API request and returns the API response. The API response is JSON decoded and stored in the value
// pointed to by v, or returned as an error if an API error has occurred. If v implements the io.Writer interface,
// the raw response will be written to v, without attempting to decode it.
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := DoRequestWithClient(ctx, c.client, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err = CheckResponse(resp); err != nil {
		return resp, err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			_, err = io.Copy(w, resp.Body)
			if err != nil {
				return nil, err
			}
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
			if err != nil {
				return nil, err
			}
		}
	}

	return resp, nil
}

// DoRequestWithClient submits an HTTP request using the specified client.
func DoRequestWithClient(ctx context.Context, client *http.Client, req *http.Request) (*http.Response, error) {
	req = req.WithContext(ctx)
	return client.Do(req)
}

// An ErrorResponse reports the error caused by an API request.
type ErrorResponse struct {
	// HTTP response that caused this error.
	Response *http.Response `json:"-"`
	// Error message.
	Message string `json:"message"`
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf(
		"%s %s: %d %s",
		r.Response.Request.Method, r.Response.Request.URL, r.Response.StatusCode, r.Message,
	)
}

// CheckResponse checks the API response for errors, and returns them if present.
// A response is considered an error if it has a status code outside the 200 range.
func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; c >= 200 && c <= 299 {
		return nil
	}

	errorResponse := &ErrorResponse{Response: r}
	data, err := ioutil.ReadAll(r.Body)
	if err == nil && len(data) > 0 {
		err := json.Unmarshal(data, errorResponse)
		if err != nil {
			errorResponse.Message = string(data)
		}
	}

	return errorResponse
}
