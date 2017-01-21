package resolvers

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// RecursiveResolver will walk the tree of the specified RootPath
// to resolve the given file name to the a path on disk. It is recommended
// to scope this as tight as possible as it is possibly quite slow
// the first time a file is requested. Once a file is found it's resolved
// path is cached to prevent further slow resolutions.
type RecursiveResolver struct {
	Path  string
	cache map[string]string
}

// Resolve will walk the tree of the specified RootPath
// to resolve the given file name to the a path on disk. It is recommended
// to scope this as tight as possible as it is possibly quite slow
// the first time a file is requested. Once a file is found it's resolved
// path is cached to prevent further slow resolutions.
func (r *RecursiveResolver) Resolve(name string) (string, error) {
	if r.cache == nil {
		r.cache = map[string]string{}
	}
	if p, ok := r.cache[name]; ok {
		return p, nil
	}
	var p string
	var err error
	var found bool
	err = filepath.Walk(r.Path, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, name) {
			found = true
			r.cache[name] = path
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

// Read will walk the tree of the specified RootPath
// to resolve the given file name to the a path on disk. It is recommended
// to scope this as tight as possible as it is possibly quite slow
// the first time a file is requested. Once a file is found it's resolved
// path is cached to prevent further slow resolutions.
func (r *RecursiveResolver) Read(name string) ([]byte, error) {
	p, err := r.Resolve(name)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadFile(p)
}
