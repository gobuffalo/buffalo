// Package meta provides minimal application metadata for Buffalo.
// This package replaces the github.com/gobuffalo/meta dependency,
// containing only the functionality that Buffalo actually uses.
package meta

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// BuildTags is a type alias for build tags used by plugins.
type BuildTags []string

// App holds metadata about the Buffalo application.
type App struct {
	Root        string
	Name        string
	Bin         string
	VCS         string
	WithPop     bool
	WithSQLite  bool
	WithWebpack bool
	WithNodeJs  bool
	WithYarn    bool
	WithDocker  bool
	WithGrifts  bool
	AsWeb       bool
	AsAPI       bool
}

// New creates App metadata for the given root path.
// If root is "." or empty, uses current working directory.
// First tries to load from config/buffalo-app.toml, then falls back to
// detecting features from the filesystem.
func New(root string) App {
	if root == "." || root == "" {
		if pwd, err := os.Getwd(); err == nil {
			root = pwd
		}
	}

	app := App{Root: root}

	// Try to load from buffalo-app.toml
	tomlPath := filepath.Join(root, "config", "buffalo-app.toml")
	if _, err := os.Stat(tomlPath); err == nil {
		// File exists, try to decode it
		if _, err := toml.DecodeFile(tomlPath, &app); err == nil {
			return app
		}
	}

	// Fall back to auto-detection (oldSchool approach)
	return autoDetect(app, root)
}

// autoDetect sets app fields by checking for files in the filesystem.
func autoDetect(app App, root string) App {
	// WithPop: check for database.yml
	if fileExists(filepath.Join(root, "database.yml")) {
		app.WithPop = true
	}

	return app
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
