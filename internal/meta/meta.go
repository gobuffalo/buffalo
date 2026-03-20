// Package meta provides application metadata for Buffalo's plugin system.
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
	Root    string `toml:"-"`
	WithPop bool   `toml:"with_pop"`
}

// New creates App metadata for the given root path.
// If root is "." or empty, uses current working directory.
// First tries to load WithPop from config/buffalo-app.toml,
// then falls back to detecting database.yml.
func New(root string) App {
	if root == "." || root == "" {
		if pwd, err := os.Getwd(); err == nil {
			root = pwd
		}
	}

	app := App{Root: root}

	tomlPath := filepath.Join(root, "config", "buffalo-app.toml")
	if _, err := os.Stat(tomlPath); err == nil {
		// TOML config exists, use it and skip auto-detection
		toml.DecodeFile(tomlPath, &app)
		return app
	}

	// No TOML config, auto-detect from filesystem
	if _, err := os.Stat(filepath.Join(root, "database.yml")); err == nil {
		app.WithPop = true
	}

	return app
}
