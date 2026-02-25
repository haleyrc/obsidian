// Package tmdb wraps the lib TMDB client, adapting its responses into movies
// domain types.
package tmdb

import (
	"context"
	"fmt"

	"github.com/haleyrc/lib/tmdb"
	"github.com/haleyrc/obsidian/internal/movies"
)

// Client wraps a lib TMDB client to implement [movies.TMDBClient].
type Client struct {
	c *tmdb.Client
}

// NewClient creates a Client authenticated with the given TMDB API access
// token.
func NewClient(accessToken string) (*Client, error) {
	c, err := tmdb.NewClient(accessToken)
	if err != nil {
		return nil, fmt.Errorf("tmdb: new client: %w", err)
	}

	return &Client{c: c}, nil
}

// GetMovie fetches movie details and credits from TMDB and returns them as a
// domain [movies.Movie].
func (c *Client) GetMovie(ctx context.Context, id int) (*movies.Movie, error) {
	detailResponse, err := c.c.GetMovieDetail(ctx, id)
	if err != nil {
		return nil, err
	}

	creditsResponse, err := c.c.GetMovieCredits(ctx, id)
	if err != nil {
		return nil, err
	}

	movie := &movies.Movie{
		Title:       detailResponse.Title,
		Synopsis:    detailResponse.Overview,
		ReleaseDate: detailResponse.ReleaseDate,
		Runtime:     detailResponse.Runtime,
		IMDBID:      detailResponse.IMDBID,
		TMDBID:      detailResponse.ID,
		PosterURL:   detailResponse.PosterURL(),
	}

	if coll := detailResponse.BelongsToCollection; coll != nil {
		movie.Collection = coll.Name
	}

	for _, genre := range detailResponse.Genres {
		movie.Genres = append(movie.Genres, genre.Name)
	}

	for _, member := range creditsResponse.Cast {
		movie.Cast = append(movie.Cast, movies.Actor{
			Name:       member.Name,
			Order:      member.Order,
			Popularity: member.Popularity,
		})
	}

	return movie, nil
}

// SearchMovie returns a list of results that match the provided title. Results
// are filtered to omit movies with a popularity < 1.0 and fewer than 100 votes
// to cut down on noise from things like behind-the-scenes featurettes, etc.
func (c *Client) SearchMovie(ctx context.Context, title string) ([]movies.MovieSearchResult, error) {
	resp, err := c.c.SearchMovie(ctx, title)
	if err != nil {
		return nil, fmt.Errorf("search movie: %w", err)
	}

	results := make([]movies.MovieSearchResult, 0, len(resp.Results))
	for _, res := range resp.Results {

		if res.Popularity < 1.0 || res.VoteCount < 100 {
			continue
		}

		results = append(results, movies.MovieSearchResult{
			ID:          res.ID,
			ReleaseDate: res.ReleaseDate,
			Synopsis:    res.Overview,
			Title:       res.Title,
		})
	}

	return results, nil
}
