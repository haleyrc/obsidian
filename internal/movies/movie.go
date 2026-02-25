package movies

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"strings"
	"text/template"
	"time"

	"github.com/haleyrc/obsidian/internal/obsidian"
	"gopkg.in/yaml.v3"
)

//go:embed movie.tmpl
var movieTemplate string

// Movie holds the raw metadata retrieved from TMDB for a single film.
type Movie struct {
	Cast        []Actor
	Collection  string
	Genres      []string
	IMDBID      string
	PosterURL   string
	ReleaseDate string
	Runtime     int
	Synopsis    string
	TMDBID      int
	Title       string
}

// MovieNote is the Obsidian-ready representation of a movie, split into YAML
// frontmatter fields and template body fields.
type MovieNote struct {
	// Frontmatter
	IMDBID      string          `yaml:"imdb_id"`
	LastWatched string          `yaml:"last_watched"`
	Owned       bool            `yaml:"owned"`
	Playlists   []obsidian.Link `yaml:"playlists"`
	TMDBID      int             `yaml:"tmdb_id"`

	// Body
	Cast        []obsidian.Link `yaml:"-"`
	Collection  obsidian.Link   `yaml:"-"`
	Genres      []obsidian.Link `yaml:"-"`
	IMDBURL     string          `yaml:"-"`
	Poster      obsidian.Link   `yaml:"-"`
	ReleaseDate string          `yaml:"-"`
	Runtime     int             `yaml:"-"`
	Synopsis    string          `yaml:"-"`
	Title       string          `yaml:"-"`
	TMDBURL     string          `yaml:"-"`
}

// Write renders the movie note as an Obsidian markdown document with YAML
// frontmatter and writes it to w.
func (mn *MovieNote) Write(w io.Writer) error {
	var buff bytes.Buffer

	fmt.Fprintln(&buff, "---")

	enc := yaml.NewEncoder(&buff)
	if err := enc.Encode(mn); err != nil {
		return fmt.Errorf("write: %w", err)
	}

	fmt.Fprintln(&buff, "---")

	funcs := template.FuncMap{
		"inline": func(links []obsidian.Link) string {
			ss := make([]string, 0, len(links))
			for _, link := range links {
				ss = append(ss, link.String())
			}
			return strings.Join(ss, ", ")
		},
		"year": func(ts string) (int, error) {
			t, err := time.Parse(time.DateOnly, ts)
			if err != nil {
				return 0, err
			}
			return t.Year(), nil
		},
	}

	t, err := template.New("movie").Funcs(funcs).Parse(movieTemplate)
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}

	if err := t.Execute(&buff, mn); err != nil {
		return err
	}

	if _, err := io.Copy(w, &buff); err != nil {
		return fmt.Errorf("write: %w", err)
	}

	return nil
}

// MovieSearchResult is a summary of a movie returned from a TMDB search query.
type MovieSearchResult struct {
	ID          int
	ReleaseDate string
	Synopsis    string
	Title       string
}
