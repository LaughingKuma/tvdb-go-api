package tvdb

import (
	"encoding/json"
	"fmt"
	"io"
)

// GetSeriesByID fetches a series by its ID
func (c *Client) GetSeriesByID(id int) (*Series, error) {
	path := fmt.Sprintf("/series/%d", id)
	
	var response struct {
		Data Series `json:"data"`
	}

	err := c.Get(path, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to get series: %w", err)
	}

	return &response.Data, nil
}


func (c *Client) GetSeriesEpisodes(seriesID int, seasonType string, page int) ([]Episode, error) {
    path := fmt.Sprintf("/series/%d/episodes/%s?page=%d", seriesID, seasonType, page)
    
    var response struct {
        Data struct {
            Episodes []Episode `json:"episodes"`
        } `json:"data"`
    }

    resp, err := c.doRequest("GET", path, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to get series episodes: %w", err)
    }
    defer resp.Body.Close()

    // Print raw response for debugging
    bodyBytes, _ := io.ReadAll(resp.Body)
    fmt.Printf("Raw response: %s\n", string(bodyBytes))

    err = json.Unmarshal(bodyBytes, &response)
    if err != nil {
        return nil, fmt.Errorf("failed to unmarshal response: %w", err)
    }

    return response.Data.Episodes, nil
}

// GetEpisodeByID fetches an episode by its ID
func (c *Client) GetEpisodeByID(id int) (*Episode, error) {
	path := fmt.Sprintf("/episodes/%d", id)
	
	var response struct {
		Data Episode `json:"data"`
	}

	err := c.Get(path, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to get episode: %w", err)
	}

	return &response.Data, nil
}

// GetSeriesSeasons fetches seasons for a series
func (c *Client) GetSeriesSeasons(seriesID int) ([]Season, error) {
	path := fmt.Sprintf("/series/%d/seasons", seriesID)
	
	var response struct {
		Data []Season `json:"data"`
	}

	err := c.Get(path, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to get series seasons: %w", err)
	}

	return response.Data, nil
}

// GetMovieByID fetches a movie by its ID
func (c *Client) GetMovieByID(id int) (*Movie, error) {
	path := fmt.Sprintf("/movies/%d", id)
	
	var response struct {
		Data Movie `json:"data"`
	}

	err := c.Get(path, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to get movie: %w", err)
	}

	return &response.Data, nil
}