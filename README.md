# tvdb-go-api

A Go client library for The TVDB API v4.

## Overview

This package provides a Go implementation for interacting with The TVDB API v4. It allows you to easily integrate TVDB data into your Go applications, providing access to information about TV series, episodes, movies, and more.

## Features

- Custom time parsing for TVDB's date-time format
- Structs for various TVDB entities:
  - Series
  - Season
  - Episode
  - Movie
  - Person
  - Artwork
  - SearchResult
- Pagination support for series episodes

## Installation

To install the package, use:

```go get github.com/LaughinKuma/tvdb-go-api```

## Usage

Check out the examples/basic_usage.go file for a demonstration of how to use the library. This example shows how to:

1. Create a TVDB client using an API key
2. Search for a TV series
3. Retrieve and print details about the series
4. Get information about the first 5 episodes of the series

```go
client := tvdb.NewClient("YOUR_API_KEY")

// Search for a TV series
results, err := client.Search("Breaking Bad")
if err != nil {
    log.Fatal(err)
}

// Print search results
for _, result := range results {
    fmt.Printf("ID: %d, Name: %s\n", result.ID, result.Name)
}

// Get details for the first search result
series, err := client.GetSeriesDetails(results[0].ID)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Overview: %s\n", series.Overview)
fmt.Printf("First Aired: %s\n", series.FirstAired)

// Get the first 5 episodes
episodes, err := client.GetSeriesEpisodes(series.ID, 1, 5)
if err != nil {
    log.Fatal(err)
}

for _, episode := range episodes {
    fmt.Printf("S%02dE%02d: %s\n", episode.SeasonNumber, episode.Number, episode.Name)
}
```

##  Structure

- `/models`: Contains the main data structures used in the API.
- `/examples`: Reserved for future usage examples.
- `/internal`: Reserved for internal package use.

## Models

The package provides several structs to represent TVDB data:

- `Series`: Represents a TV series with details like name, air dates, status, and ratings.
- `Season`: Represents a season of a TV series.
- `Episode`: Represents an individual episode of a TV series.
- `Movie`: Represents a movie with details similar to a series.
- `Person`: Represents individuals associated with series or movies.
- `Artwork`: Represents artwork associated with series, movies, or people.
- `SearchResult`: Represents a search result from the TVDB API.

## Dependencies

The project uses the following external dependencies:

- `github.com/hashicorp/go-retryablehttp`: For making HTTP requests with retry functionality.
- `github.com/stretchr/testify`: For writing and running tests.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

Enjoy it, do whatever you want with this. I just did this for fun.
If it's of service to you, go ahead. It's an example of over engineering.

## Disclaimer

This library is not officially associated with or endorsed by The TVDB.

## Credits

- [The TVDB](https://thetvdb.com/) - The source of all the TV and movie data.
- [The TVDB API Documentation](https://www.thetvdb.com/api-information) - Comprehensive guide to The TVDB API.
- [tvdb-v4-python](https://github.com/thetvdb/tvdb-v4-python) - The official Python client for The TVDB API, which served as a reference for this Go implementation.