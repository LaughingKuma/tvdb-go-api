package models

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCustomTimeUnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected time.Time
		wantErr  bool
	}{
		{
			name:     "Valid date",
			input:    `"2023-05-15 14:30:00"`,
			expected: time.Date(2023, 5, 15, 14, 30, 0, 0, time.UTC),
			wantErr:  false,
		},
		{
			name:     "Invalid date format",
			input:    `"2023-05-15"`,
			expected: time.Time{},
			wantErr:  true,
		},
		{
			name:     "Invalid JSON",
			input:    `"2023-05-15 14:30:00`,
			expected: time.Time{},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ct CustomTime
			err := json.Unmarshal([]byte(tt.input), &ct)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, time.Time(ct))
			}
		})
	}
}

func TestSeriesJSONTags(t *testing.T) {
	series := Series{
		ID:   1,
		Name: "Test Series",
	}

	jsonData, err := json.Marshal(series)
	assert.NoError(t, err)

	assert.Contains(t, string(jsonData), `"id":1`)
	assert.Contains(t, string(jsonData), `"name":"Test Series"`)
}

func TestEpisodeJSONTags(t *testing.T) {
	episode := Episode{
		ID:                 1,
		Name:               "Test Episode",
		AiredSeason:        1,
		AiredEpisodeNumber: 1,
	}

	jsonData, err := json.Marshal(episode)
	assert.NoError(t, err)

	assert.Contains(t, string(jsonData), `"id":1`)
	assert.Contains(t, string(jsonData), `"name":"Test Episode"`)
	assert.Contains(t, string(jsonData), `"airedSeason":1`)
	assert.Contains(t, string(jsonData), `"airedEpisodeNumber":1`)
}

func TestMovieJSONTags(t *testing.T) {
	movie := Movie{
		ID:   1,
		Name: "Test Movie",
	}

	jsonData, err := json.Marshal(movie)
	assert.NoError(t, err)

	assert.Contains(t, string(jsonData), `"id":1`)
	assert.Contains(t, string(jsonData), `"name":"Test Movie"`)
}

func TestSearchResultJSONTags(t *testing.T) {
	searchResult := SearchResult{
		ObjectID: "1",
		Type:     "series",
		Name:     "Test Series",
		Image:    "http://example.com/image.jpg",
		Overview: "Test overview",
	}

	jsonData, err := json.Marshal(searchResult)
	assert.NoError(t, err)

	assert.Contains(t, string(jsonData), `"objectID":"1"`)
	assert.Contains(t, string(jsonData), `"type":"series"`)
	assert.Contains(t, string(jsonData), `"name":"Test Series"`)
	assert.Contains(t, string(jsonData), `"image_url":"http://example.com/image.jpg"`)
	assert.Contains(t, string(jsonData), `"overview":"Test overview"`)
}