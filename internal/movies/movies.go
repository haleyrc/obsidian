package movies

import (
	"context"
)

// TMDBClient defines the TMDB operations required by the movies package.
type TMDBClient interface {
	GetMovie(ctx context.Context, id int) (*Movie, error)
	SearchMovie(ctx context.Context, title string) ([]MovieSearchResult, error)
}
