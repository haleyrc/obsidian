package main

import (
	"context"
	"log"
	"os"

	"github.com/haleyrc/lib/cli"
	"github.com/haleyrc/obsidian/internal/movies"
	"github.com/haleyrc/obsidian/internal/tmdb"
)

// App holds the dependencies shared across all CLI subcommands.
type App struct {
	TMDBClient movies.TMDBClient
}

func main() {
	tmdbClient, err := tmdb.NewClient(os.Getenv("TMDB_ACCESS_TOKEN"))
	if err != nil {
		log.Println("ERROR:", err)
		os.Exit(cli.ExitCodeError)
	}

	app := App{TMDBClient: tmdbClient}

	cli := cli.CLI[App]{
		Name: "movies",
		Commands: map[string]cli.Command[App]{
			"new": {
				Description: "Create a new movie note",
				Run:         runNewMovie,
			},
		},
	}

	os.Exit(cli.Run(context.Background(), app, os.Args[1:]))
}
