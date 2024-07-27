package tvdb

import (
	"github.com/LaughinKuma/tvdb-go-api/client"
	"github.com/LaughinKuma/tvdb-go-api/endpoints"
	"github.com/LaughinKuma/tvdb-go-api/models"
	"github.com/LaughinKuma/tvdb-go-api/search"
)

type TVDB struct {
	Client *client.Client
}

func New(apiKey string) (*TVDB, error) {
	c, err := client.NewClient(apiKey)
	if err != nil {
		return nil, err
	}
	return &TVDB{Client: c}, nil
}

func (t *TVDB) Search(query string) ([]models.SearchResult, error) {
	return search.Search(t.Client, query)
}

// GetSeriesByID wraps the endpoints.GetSeriesByID function
func (t *TVDB) GetSeriesByID(id int) (*models.Series, error) {
	return endpoints.GetSeriesByID(t.Client, id)
}

// GetSeriesEpisodes wraps the endpoints.GetSeriesEpisodes function
func (t *TVDB) GetSeriesEpisodes(seriesID int, seasonType string, page int) ([]models.Episode, int, int, error) {
	return endpoints.GetSeriesEpisodes(t.Client, seriesID, seasonType, page)
}

// GetEpisodeByID wraps the endpoints.GetEpisodeByID function
func (t *TVDB) GetEpisodeByID(id int) (*models.Episode, error) {
	return endpoints.GetEpisodeByID(t.Client, id)
}

// GetSeriesSeasons wraps the endpoints.GetSeriesSeasons function
func (t *TVDB) GetSeriesSeasons(seriesID int) ([]models.Season, error) {
	return endpoints.GetSeriesSeasons(t.Client, seriesID)
}

// GetMovieByID wraps the endpoints.GetMovieByID function
func (t *TVDB) GetMovieByID(id int) (*models.Movie, error) {
	return endpoints.GetMovieByID(t.Client, id)
}