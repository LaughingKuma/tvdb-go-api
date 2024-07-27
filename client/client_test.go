package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/LaughinKuma/tvdb-go-api/auth"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/stretchr/testify/assert"
)

// newTestClient creates a new Client with a custom base URL for testing
func newTestClient(apiKey, baseURL string) (*Client, error) {
	authClient := auth.NewAuth(apiKey)
	httpClient := retryablehttp.NewClient()
	httpClient.RetryMax = 3

	client := &Client{
		Auth:       authClient,
		httpClient: httpClient,
		baseURL:    baseURL,
	}

	// Mock the login process
	client.Auth.Token = "test-token"

	return client, nil
}

func TestNewClient(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "success",
			"data": map[string]string{
				"token": "test-token",
			},
		})
	}))
	defer ts.Close()

	client, err := newTestClient("test-api-key", ts.URL)

	assert.NoError(t, err)
	assert.NotNil(t, client)
	assert.Equal(t, ts.URL, client.baseURL)
	assert.Equal(t, "test-token", client.Auth.Token)
}

func TestDoRequest(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data": "test"}`))
	}))
	defer ts.Close()

	client, _ := newTestClient("test-api-key", ts.URL)

	resp, err := client.DoRequest("GET", "/test", nil)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestGet(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"data": "test"})
	}))
	defer ts.Close()

	client, _ := newTestClient("test-api-key", ts.URL)

	var result map[string]string
	err := client.Get("/test", &result)

	assert.NoError(t, err)
	assert.Equal(t, "test", result["data"])
}

func TestPost(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		var requestBody map[string]string
		json.NewDecoder(r.Body).Decode(&requestBody)
		assert.Equal(t, "test-data", requestBody["key"])
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"response": "success"})
	}))
	defer ts.Close()

	client, _ := newTestClient("test-api-key", ts.URL)

	requestBody := map[string]string{"key": "test-data"}
	var responseBody map[string]string
	err := client.Post("/test", requestBody, &responseBody)

	assert.NoError(t, err)
	assert.Equal(t, "success", responseBody["response"])
}

func TestSetBaseURL(t *testing.T) {
	client, _ := newTestClient("test-api-key", "http://old-url.com")
	client.SetBaseURL("http://new-url.com")
	assert.Equal(t, "http://new-url.com", client.baseURL)
}