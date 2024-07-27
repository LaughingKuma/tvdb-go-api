// Package endpoints provides functions for interacting with specific TVDB API endpoints.
package endpoints

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/LaughinKuma/tvdb-go-api/client"
	"github.com/LaughinKuma/tvdb-go-api/models"
)

// GetSeriesByID fetches a series by its ID.
func GetSeriesByID(c client.ClientInterface, id int) (*models.Series, error) {
	path := fmt.Sprintf("/series/%d", id)
	
	var response struct {
		Data models.Series `json:"data"`
	}

	err := c.Get(path, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to get series: %w", err)
	}

	return &response.Data, nil
}

// GetSeriesEpisodes fetches episodes for a series.
func GetSeriesEpisodes(c client.ClientInterface, seriesID int, seasonType string, page int) ([]models.Episode, int, int, error) {
    path := fmt.Sprintf("/series/%d/episodes/%s?page=%d", seriesID, seasonType, page)
    
    resp, err := c.DoRequest("GET", path, nil)
    if err != nil {
        return nil, 0, 0, fmt.Errorf("failed to get series episodes: %w", err)
    }
    defer resp.Body.Close()

    bodyBytes, _ := io.ReadAll(resp.Body)

    var response models.SeriesEpisodesResponse
    err = json.Unmarshal(bodyBytes, &response)
    if err != nil {
        return nil, 0, 0, fmt.Errorf("failed to unmarshal response: %w", err)
    }

    return response.Data.Episodes, response.Links.TotalItems, response.Links.PageSize, nil
}

// GetEpisodeByID fetches an episode by its ID.
func GetEpisodeByID(c client.ClientInterface, id int) (*models.Episode, error) {
	path := fmt.Sprintf("/episodes/%d", id)
	
	var response struct {
		Data models.Episode `json:"data"`
	}

	err := c.Get(path, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to get episode: %w", err)
	}

	return &response.Data, nil
}

// GetSeriesSeasons fetches seasons for a series.
func GetSeriesSeasons(c client.ClientInterface, seriesID int) ([]models.Season, error) {
	path := fmt.Sprintf("/series/%d/seasons", seriesID)
	
	var response struct {
		Data []models.Season `json:"data"`
	}

	err := c.Get(path, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to get series seasons: %w", err)
	}

	return response.Data, nil
}

// GetMovieByID fetches a movie by its ID.
func GetMovieByID(c client.ClientInterface, id int) (*models.Movie, error) {
	path := fmt.Sprintf("/movies/%d", id)
	
	var response struct {
		Data models.Movie `json:"data"`
	}

	err := c.Get(path, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to get movie: %w", err)
	}

	return &response.Data, nil
}