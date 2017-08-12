package packr

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// NewBox returns a Box that can be used to
// retrieve files from either disk or the embedded
// binary.
func NewBox(path string) Box {
	var cd string
	if !filepath.IsAbs(path) {
		_, filename, _, _ := runtime.Caller(1)
		cd = filepath.Dir(filename)
	}
	return Box{
		Path:       path,
		callingDir: cd,
	}
}

// Box represent a folder on a disk you want to
// have access to in the built Go binary.
type Box struct {
	Path       string
	callingDir string
	data       map[string][]byte
}

// String of the file asked for or an empty string.
func (b Box) String(name string) string {
	return string(b.Bytes(name))
}

// MustString returns either the string of the requested
// file or an error if it can not be found.
func (b Box) MustString(name string) (string, error) {
	bb, err := b.MustBytes(name)
	return string(bb), err
}

// Bytes of the file asked for or an empty byte slice.
func (b Box) Bytes(name string) []byte {
	bb, _ := b.MustBytes(name)
	return bb
}

// MustBytes returns either the byte slice of the requested
// file or an error if it can not be found.
func (b Box) MustBytes(name string) ([]byte, error) {
	f, err := b.find(name)
	if err == nil {
		bb := &bytes.Buffer{}
		bb.ReadFrom(f)
		return bb.Bytes(), err
	}
	p := filepath.Join(b.callingDir, b.Path, name)
	return ioutil.ReadFile(p)
}

func (b Box) Has(name string) bool {
	_, err := b.find(name)
	if err != nil {
		return false
	}
	return true
}

func (b Box) find(name string) (File, error) {
	name = strings.TrimPrefix(name, "/")
	name = strings.Replace(name, "\\", "/", -1)
	if _, ok := data[b.Path]; ok {
		if bb, ok := data[b.Path][name]; ok {
			return newVirtualFile(name, bb), nil
		}
	}

	p := filepath.Join(b.callingDir, b.Path, name)
	f, err := os.Open(p)
	if err != nil {
		return nil, err
	}
	return physicalFile{f}, nil
}

type WalkFunc func(string, File) error

func (b Box) Walk(wf WalkFunc) error {
	if data[b.Path] == nil {
		base := filepath.Join(b.callingDir, b.Path)
		return filepath.Walk(base, func(path string, info os.FileInfo, err error) error {
			shortPath := strings.TrimPrefix(path, base)
			if info == nil || info.IsDir() {
				return nil
			}
			f, err := os.Open(path)
			if err != nil {
				return err
			}
			return wf(shortPath, physicalFile{f})
		})
	}
	for n := range data[b.Path] {
		f, err := b.find(n)
		if err != nil {
			return err
		}
		err = wf(n, f)
		if err != nil {
			return err
		}
	}
	return nil
}

// Open returns a File using the http.File interface
func (b Box) Open(name string) (http.File, error) {
	return b.find(name)
}
