// Package movies contains domain types and vault operations for managing movie
// notes in an Obsidian vault.
package movies

// Actor represents a cast member from a movie's credits.
type Actor struct {
	Name       string
	Order      int
	Popularity float64
}
