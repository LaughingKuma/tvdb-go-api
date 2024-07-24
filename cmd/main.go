package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/LaughinKuma/tvdb-go-api/pkg/tvdb"
)

func main() {
	apiKey := os.Getenv("TVDB_API_KEY")
	if apiKey == "" {
		fmt.Println("Please set the TVDB_API_KEY environment variable")
		os.Exit(1)
	}

	client, err := tvdb.NewClient(apiKey)
	if err != nil {
		fmt.Printf("Error creating client: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("TVDB CLI")
	fmt.Println("Available commands: search, series, episodes, episode, seasons, movie, quit")

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Enter command: ")
		if !scanner.Scan() {
			break
		}

		command := scanner.Text()
		switch strings.ToLower(command) {
		case "search":
			handleSearch(client, scanner)
		case "series":
			handleGetSeries(client, scanner)
		case "episodes":
			handleGetEpisodes(client, scanner)
		case "episode":
			handleGetEpisode(client, scanner)
		case "seasons":
			handleGetSeasons(client, scanner)
		case "movie":
			handleGetMovie(client, scanner)
		case "quit":
			return
		default:
			fmt.Println("Unknown command")
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading input: %v\n", err)
	}
}

func handleSearch(client *tvdb.Client, scanner *bufio.Scanner) {
	fmt.Print("Enter search query: ")
	if !scanner.Scan() {
		return
	}
	query := scanner.Text()

	results, err := client.Search(query)
	if err != nil {
		fmt.Printf("Error performing search: %v\n", err)
		return
	}

	fmt.Printf("Found %d results:\n", len(results))
	for _, result := range results {
		fmt.Printf("- %s (%s): %s\n", result.Name, result.Type, result.Overview)
	}
}

func handleGetSeries(client *tvdb.Client, scanner *bufio.Scanner) {
	fmt.Print("Enter series ID: ")
	if !scanner.Scan() {
		return
	}
	id, err := strconv.Atoi(scanner.Text())
	if err != nil {
		fmt.Println("Invalid ID")
		return
	}

	series, err := client.GetSeriesByID(id)
	if err != nil {
		fmt.Printf("Error getting series: %v\n", err)
		return
	}

	fmt.Printf("Series: %s\n", series.Name)
	fmt.Printf("First Aired: %s\n", series.FirstAired)
	fmt.Printf("Last Updated: %s\n", time.Time(series.LastUpdated).Format("2006-01-02 15:04:05"))
	fmt.Printf("Overview: %s\n", series.Overview)
	fmt.Println("Aliases:")
	for _, alias := range series.Aliases {
		fmt.Printf("  - %s (%s)\n", alias.Name, alias.Language)
	}
}


func handleGetEpisodes(client *tvdb.Client, scanner *bufio.Scanner) {
    fmt.Print("Enter series ID: ")
    if !scanner.Scan() {
        fmt.Println("Error reading input")
        return
    }
    id, err := strconv.Atoi(scanner.Text())
    if err != nil {
        fmt.Println("Invalid ID. Please enter a numeric value.")
        return
    }

    fmt.Print("Enter page number: ")
    if !scanner.Scan() {
        fmt.Println("Error reading input")
        return
    }
    page, err := strconv.Atoi(scanner.Text())
    if err != nil {
        fmt.Println("Invalid page number. Please enter a numeric value.")
        return
    }

    fmt.Print("Enter season type (default, official, dvd, absolute, alternate, regional): ")
    if !scanner.Scan() {
        fmt.Println("Error reading input")
        return
    }
    seasonType := scanner.Text()
    if seasonType == "" {
        seasonType = "default"
    }

    fmt.Printf("Fetching episodes for series ID %d, page %d, season type '%s'...\n", id, page, seasonType)

    episodes, err := client.GetSeriesEpisodes(id, seasonType, page)
    if err != nil {
        fmt.Printf("Error getting episodes: %v\n", err)
        return
    }

    if len(episodes) == 0 {
        fmt.Println("No episodes found for this criteria.")
        return
    }

    fmt.Printf("Retrieved %d episodes:\n", len(episodes))
    for _, episode := range episodes {
        fmt.Printf("S%02dE%02d: %s (ID: %d)\n", episode.AiredSeason, episode.AiredEpisodeNumber, episode.Name, episode.ID)
    }

    fmt.Println("\nNote: If you expected more episodes, try different page numbers or season types.")
}

func handleGetEpisode(client *tvdb.Client, scanner *bufio.Scanner) {
	fmt.Print("Enter episode ID: ")
	if !scanner.Scan() {
		return
	}
	id, err := strconv.Atoi(scanner.Text())
	if err != nil {
		fmt.Println("Invalid ID")
		return
	}

	episode, err := client.GetEpisodeByID(id)
	if err != nil {
		fmt.Printf("Error getting episode: %v\n", err)
		return
	}

	fmt.Printf("Episode: %s\n", episode.Name)
	fmt.Printf("Aired Date: %s\n", time.Time(episode.AiredDate).Format("2006-01-02"))
	fmt.Printf("Last Updated: %s\n", time.Time(episode.LastUpdated).Format("2006-01-02 15:04:05"))
	fmt.Printf("Overview: %s\n", episode.Overview)
}

func handleGetSeasons(client *tvdb.Client, scanner *bufio.Scanner) {
	fmt.Print("Enter series ID: ")
	if !scanner.Scan() {
		return
	}
	id, err := strconv.Atoi(scanner.Text())
	if err != nil {
		fmt.Println("Invalid ID")
		return
	}

	seasons, err := client.GetSeriesSeasons(id)
	if err != nil {
		fmt.Printf("Error getting seasons: %v\n", err)
		return
	}

	for _, season := range seasons {
		fmt.Printf("Season %d: %s (Episodes: %d)\n", season.Number, season.Name, season.EpisodeCount)
	}
}

func handleGetMovie(client *tvdb.Client, scanner *bufio.Scanner) {
	fmt.Print("Enter movie ID: ")
	if !scanner.Scan() {
		return
	}
	id, err := strconv.Atoi(scanner.Text())
	if err != nil {
		fmt.Println("Invalid ID")
		return
	}

	movie, err := client.GetMovieByID(id)
	if err != nil {
		fmt.Printf("Error getting movie: %v\n", err)
		return
	}

	fmt.Printf("Movie: %s\nReleased: %s\nOverview: %s\n", movie.Name, movie.ReleaseDate, movie.Overview)
}