package auth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAuth(t *testing.T) {
	auth := NewAuth("test-api-key")
	assert.NotNil(t, auth)
	assert.Equal(t, "test-api-key", auth.APIKey)
	assert.Equal(t, DefaultBaseURL, auth.baseURL)
	assert.NotNil(t, auth.client)
}

func TestLogin(t *testing.T) {
	// Create a test server to mock the TVDB API
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/login", r.URL.Path)
		assert.Equal(t, "POST", r.Method)

		var requestBody map[string]string
		json.NewDecoder(r.Body).Decode(&requestBody)
		assert.Equal(t, "test-api-key", requestBody["apikey"])

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "success",
			"data": map[string]string{
				"token": "test-token",
			},
		})
	}))
	defer ts.Close()

	// Create an Auth instance with the test server URL
	auth := NewAuthWithBaseURL("test-api-key", ts.URL)

	err := auth.Login()
	assert.NoError(t, err)
	assert.Equal(t, "test-token", auth.Token)
}

func TestGetAuthHeader(t *testing.T) {
	auth := NewAuth("test-api-key")
	auth.Token = "test-token"

	header := auth.GetAuthHeader()
	assert.Equal(t, "Bearer test-token", header)
}

func TestIsAuthenticated(t *testing.T) {
	auth := NewAuth("test-api-key")
	assert.False(t, auth.IsAuthenticated())

	auth.Token = "test-token"
	assert.True(t, auth.IsAuthenticated())
}

func TestRefreshToken(t *testing.T) {
	// Create a test server to mock the TVDB API
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/login", r.URL.Path)
		assert.Equal(t, "POST", r.Method)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "success",
			"data": map[string]string{
				"token": "new-test-token",
			},
		})
	}))
	defer ts.Close()

	// Create an Auth instance with the test server URL
	auth := NewAuthWithBaseURL("test-api-key", ts.URL)
	auth.Token = "old-test-token"

	err := auth.RefreshToken()
	assert.NoError(t, err)
	assert.Equal(t, "new-test-token", auth.Token)
}