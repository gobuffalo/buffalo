package generators

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
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
	root := filepath.Join(os.Getenv("GOPATH"), "src", "github.com", "gobuffalo", "buffalo", "generators", path, "templates")
	files := Files{}
	err := filepath.Walk(root, func(p string, info os.FileInfo, err error) error {
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
