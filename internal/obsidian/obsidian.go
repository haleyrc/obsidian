// Package obsidian provides types and helpers for working with Obsidian vault
// conventions.
package obsidian

import (
	"fmt"
	"strings"
)

// SanitizeFilename replaces or removes characters that are not allowed in
// Obsidian note filenames.
func SanitizeFilename(s string) string {
	replacer := strings.NewReplacer(
		"/", "-",
		"\\", "-",
		":", " -",
		"*", "",
		"?", "",
		"\"", "",
		"<", "",
		">", "",
		"|", "-",
		".", "",
	)
	return strings.TrimSpace(replacer.Replace(s))
}

// Link is an Obsidian internal link target, rendered as [[name]] in markdown.
type Link string

// MarshalYAML implements [yaml.Marshaler], encoding the link in Obsidian's
// [[name]] syntax.
func (l Link) MarshalYAML() (any, error) {
	return l.String(), nil
}

// String returns the link formatted in Obsidian's [[name]] syntax.
func (l Link) String() string {
	s := fmt.Sprintf("[[%s]]", string(l))
	return s
}
