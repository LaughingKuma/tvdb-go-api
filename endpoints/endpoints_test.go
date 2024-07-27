// endpoints/endpoints_test.go
package endpoints

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/LaughinKuma/tvdb-go-api/client"
	"github.com/LaughinKuma/tvdb-go-api/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockClient is a mock implementation of the client.ClientInterface
type MockClient struct {
	mock.Mock
}

// Ensure MockClient implements client.ClientInterface
var _ client.ClientInterface = (*MockClient)(nil)

func (m *MockClient) Get(path string, result interface{}) error {
	args := m.Called(path, result)
	return args.Error(0)
}

func (m *MockClient) DoRequest(method, path string, body io.Reader) (*http.Response, error) {
	args := m.Called(method, path, body)
	return args.Get(0).(*http.Response), args.Error(1)
}

func (m *MockClient) Post(path string, body interface{}, result interface{}) error {
	args := m.Called(path, body, result)
	return args.Error(0)
}

func TestGetSeriesByID(t *testing.T) {
	mockClient := new(MockClient)
	seriesID := 123
	expectedSeries := &models.Series{ID: seriesID, Name: "Test Series"}

	mockClient.On("Get", "/series/123", mock.AnythingOfType("*struct { Data models.Series }")).
		Run(func(args mock.Arguments) {
			arg := args.Get(1).(*struct{ Data models.Series })
			arg.Data = *expectedSeries
		}).
		Return(nil)

	series, err := GetSeriesByID(mockClient, seriesID)

	assert.NoError(t, err)
	assert.Equal(t, expectedSeries, series)
	mockClient.AssertExpectations(t)
}

func TestGetSeriesEpisodes(t *testing.T) {
	mockClient := new(MockClient)
	seriesID := 123
	seasonType := "default"
	page := 1

	expectedResponse := models.SeriesEpisodesResponse{
		Data: struct {
			Series   models.Series    `json:"series"`
			Episodes []models.Episode `json:"episodes"`
		}{
			Episodes: []models.Episode{{ID: 1, Name: "Episode 1"}},
		},
		Links: struct {
			Prev       string `json:"prev"`
			Self       string `json:"self"`
			Next       string `json:"next"`
			TotalItems int    `json:"total_items"`
			PageSize   int    `json:"page_size"`
		}{
			TotalItems: 10,
			PageSize:   5,
		},
	}

	responseBody, _ := json.Marshal(expectedResponse)
	mockResponse := &http.Response{
		Body: io.NopCloser(bytes.NewReader(responseBody)),
	}

	mockClient.On("DoRequest", "GET", "/series/123/episodes/default?page=1", nil).
		Return(mockResponse, nil)

	episodes, totalItems, pageSize, err := GetSeriesEpisodes(mockClient, seriesID, seasonType, page)

	assert.NoError(t, err)
	assert.Equal(t, expectedResponse.Data.Episodes, episodes)
	assert.Equal(t, expectedResponse.Links.TotalItems, totalItems)
	assert.Equal(t, expectedResponse.Links.PageSize, pageSize)
	mockClient.AssertExpectations(t)
}

func TestGetEpisodeByID(t *testing.T) {
	mockClient := new(MockClient)
	episodeID := 456
	expectedEpisode := &models.Episode{ID: episodeID, Name: "Test Episode"}

	mockClient.On("Get", "/episodes/456", mock.AnythingOfType("*struct { Data models.Episode }")).
		Run(func(args mock.Arguments) {
			arg := args.Get(1).(*struct{ Data models.Episode })
			arg.Data = *expectedEpisode
		}).
		Return(nil)

	episode, err := GetEpisodeByID(mockClient, episodeID)

	assert.NoError(t, err)
	assert.Equal(t, expectedEpisode, episode)
	mockClient.AssertExpectations(t)
}

func TestGetSeriesSeasons(t *testing.T) {
	mockClient := new(MockClient)
	seriesID := 123
	expectedSeasons := []models.Season{{ID: 1, Name: "Season 1"}, {ID: 2, Name: "Season 2"}}

	mockClient.On("Get", "/series/123/seasons", mock.AnythingOfType("*struct { Data []models.Season }")).
		Run(func(args mock.Arguments) {
			arg := args.Get(1).(*struct{ Data []models.Season })
			arg.Data = expectedSeasons
		}).
		Return(nil)

	seasons, err := GetSeriesSeasons(mockClient, seriesID)

	assert.NoError(t, err)
	assert.Equal(t, expectedSeasons, seasons)
	mockClient.AssertExpectations(t)
}

func TestGetMovieByID(t *testing.T) {
	mockClient := new(MockClient)
	movieID := 789
	expectedMovie := &models.Movie{ID: movieID, Name: "Test Movie"}

	mockClient.On("Get", "/movies/789", mock.AnythingOfType("*struct { Data models.Movie }")).
		Run(func(args mock.Arguments) {
			arg := args.Get(1).(*struct{ Data models.Movie })
			arg.Data = *expectedMovie
		}).
		Return(nil)

	movie, err := GetMovieByID(mockClient, movieID)

	assert.NoError(t, err)
	assert.Equal(t, expectedMovie, movie)
	mockClient.AssertExpectations(t)
}

func TestErrorHandling(t *testing.T) {
	mockClient := new(MockClient)
	expectedError := errors.New("API error")

	// Test error handling for GetSeriesByID
	mockClient.On("Get", "/series/123", mock.Anything).Return(expectedError)
	_, err := GetSeriesByID(mockClient, 123)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get series")

	// Test error handling for GetSeriesEpisodes
	mockClient.On("DoRequest", "GET", "/series/123/episodes/default?page=1", nil).
		Return(&http.Response{}, expectedError)
	_, _, _, err = GetSeriesEpisodes(mockClient, 123, "default", 1)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get series episodes")

	// Test error handling for GetEpisodeByID
	mockClient.On("Get", "/episodes/456", mock.Anything).Return(expectedError)
	_, err = GetEpisodeByID(mockClient, 456)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get episode")

	// Test error handling for GetSeriesSeasons
	mockClient.On("Get", "/series/123/seasons", mock.Anything).Return(expectedError)
	_, err = GetSeriesSeasons(mockClient, 123)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get series seasons")

	// Test error handling for GetMovieByID
	mockClient.On("Get", "/movies/789", mock.Anything).Return(expectedError)
	_, err = GetMovieByID(mockClient, 789)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get movie")

	mockClient.AssertExpectations(t)
}