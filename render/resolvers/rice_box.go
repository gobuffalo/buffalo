package resolvers

import (
	"fmt"
	"os"
	"strings"

	rice "github.com/GeertJohan/go.rice"
)

// RiceBox uses the go.rice package to resolve files
type RiceBox struct {
	Box *rice.Box
}

// Read data from the rice.Box
func (r *RiceBox) Read(name string) ([]byte, error) {
	return r.Box.Bytes(name)
}

// Resolve the file from the rice.Box
func (r *RiceBox) Resolve(name string) (string, error) {
	var p string
	var err error
	var found bool
	err = r.Box.Walk(".", func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, name) {
			found = true
			p = path
			return err
		}
		return nil
	})
	if err != nil {
		return p, err
	}
	if !found {
		return p, fmt.Errorf("could not find file %s", name)
	}
	return p, nil
}
