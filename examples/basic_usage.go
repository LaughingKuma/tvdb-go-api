package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/LaughinKuma/tvdb-go-api"
)

func main() {
	// Replace with your actual TVDB API key
	apiKey := "your_api_key_here"

	// Create a new TVDB client
	tvdbClient, err := tvdb.New(apiKey)
	if err != nil {
		log.Fatalf("Failed to create TVDB client: %v", err)
	}

// Search for a TV series
searchQuery := "Breaking Bad"
searchResults, err := tvdbClient.Search(searchQuery)
if err != nil {
	log.Fatalf("Failed to search: %v", err)
}

fmt.Printf("Search results for '%s':\n", searchQuery)
for _, result := range searchResults {
	fmt.Printf("- %s (ID: %s)\n", result.Name, result.ObjectID)
}

// Get details for a specific series (using the first search result)
if len(searchResults) > 0 {
	// Convert the ObjectID (string) to an integer
	seriesID, err := strconv.Atoi(searchResults[0].ObjectID)
	if err != nil {
		log.Fatalf("Failed to convert series ID to integer: %v", err)
	}

	series, err := tvdbClient.GetSeriesByID(seriesID)
	if err != nil {
		log.Fatalf("Failed to get series details: %v", err)
	}

	fmt.Printf("\nSeries details for '%s':\n", series.Name)
	fmt.Printf("Overview: %s\n", series.Overview)
	fmt.Printf("First Aired: %s\n", series.FirstAired)

	// Get episodes for the series
	episodes, _, _, err := tvdbClient.GetSeriesEpisodes(seriesID, "default", 1)
	if err != nil {
		log.Fatalf("Failed to get series episodes: %v", err)
	}

	fmt.Printf("\nFirst 5 episodes:\n")
	for i, episode := range episodes {
		if i >= 5 {
			break
		}
		fmt.Printf("- S%02dE%02d: %s\n", episode.AiredSeason, episode.AiredEpisodeNumber, episode.Name)
	}
}
}