package openweather

import (
	"errors"
	"net/http"
	"net/url"
)

// Opt is a configuration option to initialize a client.
type Opt func(*Client) error

// WithHTTPClient sets the HTTP client which will be used to make requests.
func WithHTTPClient(httpClient *http.Client) Opt {
	return func(c *Client) error {
		if httpClient == nil {
			return errors.New("*http.Client: cannot be nil")
		}
		c.client = httpClient
		return nil
	}
}

// WithUserAgent sets the User-Agent header for requests made with the client.
func WithUserAgent(ua string) Opt {
	return func(c *Client) error {
		c.UserAgent = ua
		return nil
	}
}

// WithBaseURL sets the base URL for the client to make requests to.
func WithBaseURL(u string) Opt {
	return func(c *Client) error {
		url, err := url.Parse(u)
		if err != nil {
			return err
		}
		c.BaseURL = url
		return nil
	}
}
