// Package env provides environment variable utilities for Buffalo.
// This package replaces the github.com/gobuffalo/envy dependency.
package env

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
)

// Load loads .env file(s) into environment.
// If no filenames are provided, it loads the .env file in the current directory.
func Load(filenames ...string) error {
	return godotenv.Load(filenames...)
}

// Get returns the value of the environment variable named by key.
// If the variable is not set or is empty, it returns the fallback value.
func Get(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// MustGet returns the value of the environment variable named by key.
// It returns an error if the variable is not set or is empty.
func MustGet(key string) (string, error) {
	v := os.Getenv(key)
	if v == "" {
		return "", fmt.Errorf("environment variable %s is not set", key)
	}
	return v, nil
}

// Environ returns a copy of the environment variables.
func Environ() []string {
	return os.Environ()
}

// GoPath returns the GOPATH environment variable.
// If GOPATH is not set, it returns the default GOPATH: $HOME/go
func GoPath() string {
	if gp := os.Getenv("GOPATH"); gp != "" {
		return gp
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, "go")
}

// Set sets the environment variable named by key to value.
func Set(key, value string) {
	os.Setenv(key, value)
}

// Temp executes f with temporarily modified environment variables.
// After f returns, the environment is restored to its original state.
func Temp(f func()) {
	oldEnv := os.Environ()
	defer func() {
		os.Clearenv()
		for _, e := range oldEnv {
			if i := strings.IndexByte(e, '='); i > 0 {
				os.Setenv(e[:i], e[i+1:])
			}
		}
	}()
	os.Clearenv()
	f()
}
