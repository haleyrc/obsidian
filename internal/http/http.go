// Package http provides HTTP utility functions.
package http

import (
	"fmt"
	"io"
	"net/http"
)

// Download fetches the contents of url and writes them to w.
func Download(w io.Writer, url string) error {
	response, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("download: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("download: bad status: %s", response.Status)
	}

	_, err = io.Copy(w, response.Body)
	if err != nil {
		return fmt.Errorf("download: %w", err)
	}

	return nil
}
