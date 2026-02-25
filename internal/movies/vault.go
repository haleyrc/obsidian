package movies

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/haleyrc/obsidian/internal/http"
	"github.com/haleyrc/obsidian/internal/obsidian"
)

// Vault is the primary entrypoint for reading from and writing to the movies
// vault filesystem.
type Vault struct {
	Dir            string
	AttachmentsDir string
	MoviesDir      string

	ActorsDir   string
	ActorsCache []string

	GenresDir   string
	GenresCache []string
}

// LoadVault initializes a Vault rooted at dir, populating the actor and genre
// caches from the existing vault contents.
func LoadVault(dir string) (*Vault, error) {
	vault := &Vault{
		Dir:            dir,
		ActorsDir:      filepath.Join(dir, "actors"),
		AttachmentsDir: filepath.Join(dir, "attachments"),
		GenresDir:      filepath.Join(dir, "genres"),
		MoviesDir:      filepath.Join(dir, "movies"),
	}

	if err := vault.loadActors(); err != nil {
		return nil, fmt.Errorf("load vault: %w", err)
	}

	if err := vault.loadGenres(); err != nil {
		return nil, fmt.Errorf("load vault: %w", err)
	}

	return vault, nil
}

// CreateActor creates an empty markdown note for the named actor in the vault's
// actors directory.
func (vault *Vault) CreateActor(name string) error {
	path := filepath.Join(vault.ActorsDir, name+".md")
	log.Println("Creating actor file:", path)

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return nil
}

// CreateGenre creates an empty markdown note for the named genre in the vault's
// genres directory.
func (vault *Vault) CreateGenre(name string) error {
	path := filepath.Join(vault.GenresDir, name+".md")
	log.Println("Creating genre file:", path)

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return nil
}

// CreateMovie writes a complete movie note to the vault, including the poster
// image and any new actor and genre notes.
func (vault *Vault) CreateMovie(movie *Movie) error {
	title := obsidian.SanitizeFilename(movie.Title)
	path := filepath.Join(vault.MoviesDir, title+".md")
	posterFilename := title + ".jpg"

	if err := vault.downloadPoster(posterFilename, movie); err != nil {
		return fmt.Errorf("create movie: %w", err)
	}

	log.Println("Creating movie file:", path)

	note := &MovieNote{
		// Frontmatter
		IMDBID:    movie.IMDBID,
		Owned:     false,
		Playlists: []obsidian.Link{"Newly Added"},
		TMDBID:    movie.TMDBID,

		// Body
		Collection:  obsidian.Link(obsidian.SanitizeFilename(movie.Collection)),
		IMDBURL:     fmt.Sprintf("https://m.imdb.com/title/%s", movie.IMDBID),
		Poster:      obsidian.Link(posterFilename),
		ReleaseDate: movie.ReleaseDate,
		Runtime:     movie.Runtime,
		Synopsis:    movie.Synopsis,
		Title:       title,
		TMDBURL:     fmt.Sprintf("https://www.themoviedb.org/movie/%d", movie.TMDBID),
	}

	for _, actor := range filterCast(movie.Cast, vault.ActorsCache) {
		note.Cast = append(note.Cast, obsidian.Link(actor))
		if err := vault.CreateActor(actor); err != nil {
			return fmt.Errorf("create movie: %w", err)
		}
	}

	for _, genre := range movie.Genres {
		note.Genres = append(note.Genres, obsidian.Link(genre))
		if err := vault.CreateGenre(genre); err != nil {
			return fmt.Errorf("create movie: %w", err)
		}
	}

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create movie: %w", err)
	}
	defer f.Close()

	if err := note.Write(f); err != nil {
		return fmt.Errorf("create movie: %w", err)
	}

	return nil
}

func (vault *Vault) downloadPoster(filename string, movie *Movie) error {
	log.Println("Downloading poster:", movie.PosterURL)

	var buff bytes.Buffer
	if err := http.Download(&buff, movie.PosterURL); err != nil {
		return fmt.Errorf("download poster: %w", err)
	}

	path := filepath.Join(vault.AttachmentsDir, filename)
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("download poster: %w", err)
	}
	defer f.Close()

	if _, err := io.Copy(f, &buff); err != nil {
		return fmt.Errorf("download poster: %w", err)
	}

	return nil
}

func (vault *Vault) loadActors() error {
	log.Println("Loading actors from:", vault.ActorsDir)

	actors, err := loadFilenames(vault.ActorsDir)
	if err != nil {
		return fmt.Errorf("load actors: %w", err)
	}
	log.Printf("  Loaded %d actors.", len(actors))

	vault.ActorsCache = actors

	return nil
}

func (vault *Vault) loadGenres() error {
	log.Println("Loading genres from:", vault.GenresDir)

	genres, err := loadFilenames(vault.GenresDir)
	if err != nil {
		return fmt.Errorf("load genres: %w", err)
	}
	log.Printf("  Loaded %d genres.", len(genres))

	vault.GenresCache = genres

	return nil
}

func loadFilenames(dir string) ([]string, error) {
	files, err := filepath.Glob(filepath.Join(dir, "*.md"))
	if err != nil {
		return nil, err
	}

	names := make([]string, 0, len(files))
	for _, f := range files {
		names = append(names, strings.TrimSuffix(filepath.Base(f), ".md"))
	}

	return names, nil
}

func filterCast(cast []Actor, knownActors []string) []string {
	names := []string{}

	log.Println("Adding cast members")
	skipped := 0
	for _, star := range cast {
		if star.Order < 10 {
			log.Printf("  %s (Reason: order=%d)", star.Name, star.Order)
			names = append(names, star.Name)
			continue
		}

		if star.Popularity >= 1.0 {
			log.Printf("  %s (Reason: popularity=%f)", star.Name, star.Popularity)
			names = append(names, star.Name)
			continue
		}

		if slices.Contains(knownActors, star.Name) {
			log.Printf("  %s (Reason: already known)", star.Name)
			names = append(names, star.Name)
			continue
		}

		skipped += 1
	}
	log.Printf("Skipped %d cast members", skipped)

	return names
}
