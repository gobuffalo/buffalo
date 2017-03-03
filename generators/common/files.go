package common

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type File struct {
	ReadPath  string
	WritePath string
	Body      string
}

type Files []File

func Find(path string) (Files, error) {
	root := filepath.Join(os.Getenv("GOPATH"), "src", "github.com", "gobuffalo", "buffalo", "generators", path)
	files := Files{}
	err := filepath.Walk(root, func(p string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			if filepath.Ext(p) == ".tmpl" {
				f := File{ReadPath: p}
				base := filepath.Base(p)
				base = strings.TrimSuffix(base, ".tmpl")
				f.WritePath = filepath.Join(strings.Split(base, "-")...)
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
