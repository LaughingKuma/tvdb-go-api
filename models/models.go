package models

import (
	"time"
)

// CustomTime is a custom time type to handle TVDB's date-time format
type CustomTime time.Time

type Alias struct {
	Language string `json:"language"`
	Name     string `json:"name"`
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (ct *CustomTime) UnmarshalJSON(b []byte) error {
	s := string(b)
	// Remove quotes
	s = s[1 : len(s)-1]
	
	// Parse the time using the format returned by the API
	t, err := time.Parse("2006-01-02 15:04:05", s)
	if err != nil {
		return err
	}
	
	*ct = CustomTime(t)
	return nil
}

// Series represents a TV series
type Series struct {
	ID                   int         `json:"id"`
	Name                 string      `json:"name"`
	Slug                 string      `json:"slug"`
	Image                string      `json:"image"`
	FirstAired           string      `json:"firstAired"`
	LastAired            string      `json:"lastAired"`
	NextAired            string      `json:"nextAired"`
	Status               Status      `json:"status"`
	Overview             string      `json:"overview"`
	Network              string      `json:"network"`
	Runtime              int         `json:"runtime"`
	Language             string      `json:"language"`
	Genre                []string    `json:"genre"`
	LastUpdated          CustomTime  `json:"lastUpdated"`
	AverageRating        float64     `json:"averageRating"`
	OriginalCountry      string      `json:"originalCountry"`
	OriginalLanguage     string      `json:"originalLanguage"`
	ContentRating        string      `json:"contentRating"`
	ImdbID               string      `json:"imdbId"`
	ZapID                string      `json:"zap2itId"`
	Aliases              []Alias     `json:"aliases"`
	NameTranslations     []string    `json:"nameTranslations"`
	OverviewTranslations []string    `json:"overviewTranslations"`
}


type Season struct {
	ID                 int       `json:"id"`
	SeriesID           int       `json:"seriesId"`
	Number             int       `json:"number"`
	Name               string    `json:"name"`
	EpisodeCount       int       `json:"episodeCount"`
	Overview           string    `json:"overview"`
	Image              string    `json:"image"`
	NetworkID          int       `json:"networkId"`
	LastUpdated        time.Time `json:"lastUpdated"`
	NameTranslations   []string  `json:"nameTranslations"`
	OverviewTranslations []string  `json:"overviewTranslations"`
}

// Episode represents a TV episode
type Episode struct {
	ID                 int         `json:"id"`
	SeriesID           int         `json:"seriesId"`
	Name               string      `json:"name"`
	AiredSeason        int         `json:"airedSeason"`
	AiredEpisodeNumber int         `json:"airedEpisodeNumber"`
	AiredDate          CustomTime  `json:"airedDate"`
	Runtime            int         `json:"runtime"`
	Overview           string      `json:"overview"`
	Image              string      `json:"image"`
	ImdbID             string      `json:"imdbId"`
	LastUpdated        CustomTime  `json:"lastUpdated"`
	NameTranslations   []string    `json:"nameTranslations"`
	OverviewTranslations []string  `json:"overviewTranslations"`
}


// Movie represents a movie
type Movie struct {
	ID                   int         `json:"id"`
	Name                 string      `json:"name"`
	Slug                 string      `json:"slug"`
	Image                string      `json:"image"`
	ReleaseDate          CustomTime  `json:"releaseDate"`
	Status               Status      `json:"status"`
	Overview             string      `json:"overview"`
	Runtime              int         `json:"runtime"`
	Language             string      `json:"language"`
	Genre                []string    `json:"genre"`
	LastUpdated          CustomTime  `json:"lastUpdated"`
	AverageRating        float64     `json:"averageRating"`
	OriginalCountry      string      `json:"originalCountry"`
	OriginalLanguage     string      `json:"originalLanguage"`
	ContentRating        string      `json:"contentRating"`
	ImdbID               string      `json:"imdbId"`
	Aliases              []Alias     `json:"aliases"`
	NameTranslations     []string    `json:"nameTranslations"`
	OverviewTranslations []string    `json:"overviewTranslations"`
}

// Status represents the status of a series or movie
type Status struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Person represents an actor, director, or other person associated with a series or movie
type Person struct {
	ID          int         `json:"id"`
	Name        string      `json:"name"`
	Image       string      `json:"image"`
	BirthDate   CustomTime  `json:"birthDate"`
	DeathDate   CustomTime  `json:"deathDate"`
	Gender      int         `json:"gender"`
	LastUpdated CustomTime  `json:"lastUpdated"`
}

// Artwork represents artwork associated with a series, movie, or person
type Artwork struct {
	ID       int    `json:"id"`
	Language string `json:"language"`
	Type     int    `json:"type"`
	Score    int    `json:"score"`
	URL      string `json:"url"`
	Thumbnail string `json:"thumbnail"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
}

// SearchResult represents a search result from the TVDB API
type SearchResult struct {
	ObjectID string `json:"objectID"`
	Type     string `json:"type"`
	Name     string `json:"name"`
	Image    string `json:"image_url"`
	Overview string `json:"overview"`
}

type SeriesEpisodesResponse struct {
	Status string `json:"status"`
	Data   struct {
		Series   Series    `json:"series"`
		Episodes []Episode `json:"episodes"`
	} `json:"data"`
	Links struct {
		Prev       string `json:"prev"`
		Self       string `json:"self"`
		Next       string `json:"next"`
		TotalItems int    `json:"total_items"`
		PageSize   int    `json:"page_size"`
	} `json:"links"`
}
