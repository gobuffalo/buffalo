// Package meta provides minimal application metadata for Buffalo.
// This package replaces the github.com/gobuffalo/meta dependency,
// containing only the functionality that Buffalo actually uses.
package meta

import (
	"os"
	"path/filepath"
)

// BuildTags is a type alias for build tags used by plugins.
type BuildTags []string

// App holds minimal metadata about the Buffalo application.
type App struct {
	Root    string // Project root directory
	WithPop bool   // Has database.yml (uses Pop ORM)
}

// New creates App metadata for the given root path.
// If root is "." or empty, uses current working directory.
func New(root string) App {
	if root == "." || root == "" {
		if pwd, err := os.Getwd(); err == nil {
			root = pwd
		}
	}

	fileExists := func(path string) bool {
		_, err := os.Stat(path)
		return err == nil
	}

	return App{
		Root:    root,
		WithPop: fileExists(filepath.Join(root, "database.yml")),
	}
}
