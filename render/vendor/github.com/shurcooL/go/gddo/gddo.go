// Package gddo is a simple client library for accessing the godoc.org API.
//
// It provides a single utility to fetch the importers of a Go package.
package gddo

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Client manages communication with the godoc.org API.
type Client struct {
	// UserAgent is used for outbound requests to godoc.org API, if set to non-empty value.
	UserAgent string
}

// GetImporters fetches the importers of Go package with specified importPath via godoc.org API.
func (c *Client) GetImporters(importPath string) (Importers, error) {
	req, err := http.NewRequest("GET", "https://api.godoc.org/importers/"+importPath, nil)
	if err != nil {
		return Importers{}, err
	}
	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return Importers{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return Importers{}, fmt.Errorf("non-200 status code: %v", resp.StatusCode)
	}
	var importers Importers
	err = json.NewDecoder(resp.Body).Decode(&importers)
	if err != nil {
		return Importers{}, err
	}
	return importers, nil
}

// Importers contains the list of Go packages that import a given Go package.
type Importers struct {
	Results []Package
}

// Package represents a Go package.
type Package struct {
	Path     string // Import path of the package.
	Synopsis string // Synopsis of the package.
}
