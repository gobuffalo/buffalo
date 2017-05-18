package generators

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/envy"
)

// File represents the file to be templated
type File struct {
	ReadPath  string
	WritePath string
	Body      string
}

// Files is a slice of File
type Files []File

// Find all the .tmpl files inside the buffalo GOPATH
func Find(path string) (Files, error) {
	gp, err := goPath()
	if err != nil {
		return nil, err
	}
	root := filepath.Join(gp, "src", "github.com", "gobuffalo", "buffalo", "generators", path, "templates")
	files := Files{}
	err = filepath.Walk(root, func(p string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			if filepath.Ext(p) == ".tmpl" {
				f := File{ReadPath: p}
				rel := strings.TrimPrefix(p, root)

				paths := strings.Split(rel, string(os.PathSeparator))

				li := len(paths) - 1
				base := paths[li]
				base = strings.TrimSuffix(base, ".tmpl")
				if strings.HasPrefix(base, "dot-") {
					base = "." + strings.TrimPrefix(base, "dot-")
				}
				paths[li] = base
				f.WritePath = filepath.Join(paths...)

				b, err := ioutil.ReadFile(p)
				if err != nil {
					return err
				}
				f.Body = string(b)
				files = append(files, f)
			}
		}
		return nil
	})
	return files, err
}

func goPath() (string, error) {
	for _, path := range envy.GoPaths() {
		pp := filepath.Join(path, "src", "github.com", "gobuffalo", "buffalo")
		if exists(pp) {
			return path, nil
		}
	}
	return "", errors.New("buffalo was not found")
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
