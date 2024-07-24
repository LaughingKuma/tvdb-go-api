package tvdb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

// Client represents the TVDB API client
type Client struct {
	*Auth
	httpClient *retryablehttp.Client
	baseURL    string
}

// NewClient creates a new TVDB API client
func NewClient(apiKey string) (*Client, error) {
	auth := NewAuth(apiKey)
	httpClient := retryablehttp.NewClient()
	httpClient.RetryMax = 3
	httpClient.RetryWaitMin = 1 * time.Second
	httpClient.RetryWaitMax = 5 * time.Second

	client := &Client{
		Auth:       auth,
		httpClient: httpClient,
		baseURL:    baseURL,
	}

	err := client.Auth.Login()
	if err != nil {
		return nil, fmt.Errorf("failed to authenticate: %w", err)
	}

	return client, nil
}

// doRequest performs an HTTP request and handles authentication
func (c *Client) doRequest(method, path string, body io.Reader) (*http.Response, error) {
	url := c.baseURL + path
	req, err := retryablehttp.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(authHeader, c.GetAuthHeader())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}

	if resp.StatusCode == http.StatusUnauthorized {
		// Token might be expired, try to refresh
		err = c.Auth.RefreshToken()
		if err != nil {
			return nil, fmt.Errorf("error refreshing token: %w", err)
		}

		// Retry the request with the new token
		req.Header.Set(authHeader, c.GetAuthHeader())
		resp, err = c.httpClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("error sending request after token refresh: %w", err)
		}
	}

	return resp, nil
}

// Get performs a GET request to the specified path
func (c *Client) Get(path string, result interface{}) error {
	resp, err := c.doRequest("GET", path, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	err = json.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		return fmt.Errorf("error decoding response: %w", err)
	}

	return nil
}

// Post performs a POST request to the specified path
func (c *Client) Post(path string, body interface{}, result interface{}) error {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("error marshaling request body: %w", err)
	}

	resp, err := c.doRequest("POST", path, bytes.NewReader(jsonBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	err = json.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		return fmt.Errorf("error decoding response: %w", err)
	}

	return nil
}

// SetBaseURL allows changing the base URL for API requests
func (c *Client) SetBaseURL(url string) {
	c.baseURL = url
}