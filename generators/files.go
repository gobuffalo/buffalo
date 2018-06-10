package generators

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/packr"
)

// File represents the file to be templated
type File struct {
	ReadPath  string
	WritePath string
	Body      string
}

// Files is a slice of File
type Files []File

// FindByBox all the .tmpl files inside the packr.Box
func FindByBox(box packr.Box) (Files, error) {
	files := Files{}
	err := box.Walk(func(p string, file packr.File) error {
		if filepath.Ext(p) == ".tmpl" {
			p = strings.TrimPrefix(p, "/")
			f := File{ReadPath: p}
			p = strings.Replace(p, "dot-", ".", 1)
			p = strings.Replace(p, ".tmpl", "", 1)
			f.WritePath = p
			b, err := ioutil.ReadAll(file)
			if err != nil {
				return err
			}
			f.Body = string(b)
			files = append(files, f)
		}
		return nil
	})
	return files, err
}

// TemplatesPath is the "base" path for generator templates
var TemplatesPath = filepath.Join("github.com", "gobuffalo", "buffalo", "generators")
