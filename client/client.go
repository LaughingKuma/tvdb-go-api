// Package client provides the main client for interacting with the TVDB API.
package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/LaughinKuma/tvdb-go-api/auth"
	"github.com/hashicorp/go-retryablehttp"
)

// Client represents the TVDB API client

type ClientInterface interface {
    Get(path string, result interface{}) error
    DoRequest(method, path string, body io.Reader) (*http.Response, error)
    Post(path string, body interface{}, result interface{}) error
    SetBaseURL(url string)
}

type Client struct {
	Auth       *auth.Auth
	httpClient *retryablehttp.Client
	baseURL    string
}

// NewClient creates a new TVDB API client
func NewClient(apiKey string) (*Client, error) {
	authClient := auth.NewAuth(apiKey)
	httpClient := retryablehttp.NewClient()
	httpClient.RetryMax = 3
	httpClient.RetryWaitMin = 1 * time.Second
	httpClient.RetryWaitMax = 5 * time.Second

	client := &Client{
		Auth:       authClient,
		httpClient: httpClient,
		baseURL:    auth.DefaultBaseURL,
	}

	err := client.Auth.Login()
	if err != nil {
		return nil, fmt.Errorf("failed to authenticate: %w", err)
	}

	return client, nil
}

// DoRequest performs an HTTP request and handles authentication
func (c *Client) DoRequest(method, path string, body io.Reader) (*http.Response, error) {
	url := c.baseURL + path
	req, err := retryablehttp.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(auth.AuthHeader, c.Auth.GetAuthHeader())

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
		req.Header.Set(auth.AuthHeader, c.Auth.GetAuthHeader())
		resp, err = c.httpClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("error sending request after token refresh: %w", err)
		}
	}

	return resp, nil
}

// Get performs a GET request to the specified path
func (c *Client) Get(path string, result interface{}) error {
	resp, err := c.DoRequest("GET", path, nil)
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

	resp, err := c.DoRequest("POST", path, bytes.NewReader(jsonBody))
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