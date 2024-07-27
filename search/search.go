package search

import (
	"fmt"
	"net/url"

	"github.com/LaughinKuma/tvdb-go-api/models"
)

// ClientInterface defines the methods we need from the client
type ClientInterface interface {
	Get(path string, result interface{}) error
}

func Search(c ClientInterface, query string) ([]models.SearchResult, error) {
	path := fmt.Sprintf("/search?query=%s", url.QueryEscape(query))
	
	var response struct {
		Data []models.SearchResult `json:"data"`
	}

	err := c.Get(path, &response)
	if err != nil {
		return nil, fmt.Errorf("search request failed: %w", err)
	}

	return response.Data, nil
}