package tvdb

import (
	"github.com/LaughinKuma/tvdb-go-api/client"
	"github.com/LaughinKuma/tvdb-go-api/models"
	"github.com/LaughinKuma/tvdb-go-api/search"
)

type TVDB struct {
	client *client.Client
}

func New(apiKey string) (*TVDB, error) {
	c, err := client.NewClient(apiKey)
	if err != nil {
		return nil, err
	}
	return &TVDB{client: c}, nil
}

func (t *TVDB) Search(query string) ([]models.SearchResult, error) {
	return search.Search(t.client, query)
}
