package resolvers

import (
	"fmt"
	"os"
	"strings"
	"sync"

	rice "github.com/GeertJohan/go.rice"
)

var moot = &sync.Mutex{}

// RiceBox uses the go.rice package to resolve files
type RiceBox struct {
	Path  string
	found bool
	box   *rice.Box
}

func (r *RiceBox) findBox() error {
	moot.Lock()
	defer moot.Unlock()
	var err error
	if !r.found {
		r.box, err = rice.FindBox(r.Path)
		if err != nil {
			return err
		}
		r.found = true
	}
	return nil
}

// Read data from the rice.Box
func (r *RiceBox) Read(name string) ([]byte, error) {
	err := r.findBox()
	if err != nil {
		return nil, err
	}
	return r.box.Bytes(name)
}

// Resolve the file from the rice.Box
func (r *RiceBox) Resolve(name string) (string, error) {
	err := r.findBox()
	if err != nil {
		return "", err
	}
	var p string
	var found bool
	err = r.box.Walk(".", func(path string, info os.FileInfo, err error) error {
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
