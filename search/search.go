package tvdb

import (
	"fmt"
	"net/url"
)

func (c *Client) Search(query string) ([]SearchResult, error) {
	path := fmt.Sprintf("/search?query=%s", url.QueryEscape(query))
	
	var response struct {
		Data []SearchResult `json:"data"`
	}

	err := c.Get(path, &response)
	if err != nil {
		return nil, fmt.Errorf("search request failed: %w", err)
	}

	return response.Data, nil
}