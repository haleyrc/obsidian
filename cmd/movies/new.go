package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/haleyrc/lib/tmdb"
	"github.com/haleyrc/obsidian/internal/movies"
)

func runNewMovie(ctx context.Context, app App, args []string) error {
	fs := flag.NewFlagSet("newmovie", flag.ContinueOnError)

	var (
		dir   = fs.String("dir", ".", "The root of the vault")
		id    = fs.Int("id", 0, "The TMDB ID of the movie; overrides title")
		title = fs.String("title", "", "The title of the movie to add")
	)

	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("new movie: %w", err)
	}

	if *id == 0 && *title == "" {
		return fmt.Errorf("new movie: you must specify either an id or a title")
	}

	movie, err := getOrSearchMovie(ctx, app.TMDBClient, *id, *title)
	if err != nil {
		return fmt.Errorf("new movie: %w", err)
	}

	vault, err := movies.LoadVault(*dir)
	if err != nil {
		return fmt.Errorf("new movie: %w", err)
	}

	if err := vault.CreateMovie(movie); err != nil {
		return fmt.Errorf("new movie: %w", err)
	}

	return nil
}

func getMovie(ctx context.Context, c movies.TMDBClient, id int) (*movies.Movie, error) {
	log.Println("Getting movie:", id)

	movie, err := c.GetMovie(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get movie: %w", err)
	}

	log.Println("Summary:")
	log.Println("  Synopsis:", truncate(40, movie.Synopsis))
	log.Println("  Release date:", movie.ReleaseDate)
	log.Printf("  Runtime: %d minutes", movie.Runtime)
	log.Println("  IMDB ID:", movie.IMDBID)
	log.Println("  TMDB ID:", movie.TMDBID)
	log.Println("  Collection:", movie.Collection)
	log.Println("  Cast Members:", len(movie.Cast))
	log.Println("  Genres:", len(movie.Genres))
	log.Println("  Poster URL:", movie.PosterURL)

	return movie, nil
}

func getOrSearchMovie(ctx context.Context, c movies.TMDBClient, id int, title string) (*movies.Movie, error) {
	if id == 0 {
		return searchMovie(ctx, c, title)
	}
	return getMovie(ctx, c, id)
}

func searchMovie(ctx context.Context, c movies.TMDBClient, title string) (*movies.Movie, error) {
	log.Println("Searching for movie:", title)

	results, err := c.SearchMovie(ctx, title)
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no results for movie: %s", title)
	}

	if len(results) > 1 {
		log.Printf("Found %d results for %q:", len(results), title)
		log.Println()
		for _, result := range results {
			log.Printf("=== %s (%s)", result.Title, result.ReleaseDate)
			log.Printf("    %s", tmdb.MovieURL(result.ID))
			log.Printf("    %s", truncate(120, result.Synopsis))
			log.Println()
		}
		log.Println("Full search results:", tmdb.SearchURL(title))
		log.Println()
		return nil, fmt.Errorf("multiple results found: run again with the --id flag")
	}

	return getMovie(ctx, c, results[0].ID)
}

func truncate(n int, s string) string {
	if len(s) < n {
		return s
	}
	return s[0:n-3] + "..."
}
