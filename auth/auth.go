// Package auth provides authentication functionality for the TVDB API.
package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

const (
	// BaseURL is the base URL for the TVDB API.
	DefaultBaseURL = "https://api4.thetvdb.com/v4"
	
	loginPath = "/login"
	
	// AuthHeader is the name of the authorization header used in API requests.
	AuthHeader = "Authorization"
)

// Auth holds the authentication information for the TVDB API.
type Auth struct {
	APIKey  string
	Token   string
	client  *retryablehttp.Client
	baseURL string
}


type loginResponse struct {
	Status string `json:"status"`
	Data   struct {
		Token string `json:"token"`
	} `json:"data"`
}

// NewAuth creates a new Auth instance.
func NewAuth(apiKey string) *Auth {
	return NewAuthWithBaseURL(apiKey, DefaultBaseURL)
}


func NewAuthWithBaseURL(apiKey, baseURL string) *Auth {
	client := retryablehttp.NewClient()
	client.RetryMax = 3
	client.RetryWaitMin = 1 * time.Second
	client.RetryWaitMax = 5 * time.Second

	return &Auth{
		APIKey:  apiKey,
		client:  client,
		baseURL: baseURL,
	}
}

// Login authenticates with the TVDB API and obtains a token.
func (a *Auth) Login() error {
	url := a.baseURL + loginPath

	body := map[string]string{"apikey": a.APIKey}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("error marshaling login request: %w", err)
	}

	req, err := retryablehttp.NewRequest("POST", url, jsonBody)
	if err != nil {
		return fmt.Errorf("error creating login request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending login request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("login failed with status: %s", resp.Status)
	}

	var loginResp loginResponse
	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
		return fmt.Errorf("error decoding login response: %w", err)
	}

	if loginResp.Data.Token == "" {
		return fmt.Errorf("no token received in login response")
	}

	a.Token = loginResp.Data.Token
	return nil
}

// GetAuthHeader returns the authorization header for API requests.
func (a *Auth) GetAuthHeader() string {
	return fmt.Sprintf("Bearer %s", a.Token)
}

// IsAuthenticated checks if the current token is valid.
func (a *Auth) IsAuthenticated() bool {
	return a.Token != ""
}

// RefreshToken attempts to refresh the authentication token.
func (a *Auth) RefreshToken() error {
	// For TVDB API v4, we simply re-login to refresh the token
	return a.Login()
}