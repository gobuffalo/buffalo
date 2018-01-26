package generators

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/packr"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
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

// Find all the .tmpl files inside the buffalo GOPATH
func Find(path string) (Files, error) {
	warningMsg := "Find is deprecated, and will be removed in v0.12.0. Use generators.FindByBox instead."
	_, file, no, ok := runtime.Caller(1)
	if ok {
		warningMsg = fmt.Sprintf("%s Called from %s:%d", warningMsg, file, no)
	}

	logrus.Info(warningMsg)
	mu := &sync.Mutex{}
	wg := &errgroup.Group{}
	files := Files{}
	for _, gp := range envy.GoPaths() {
		func(gp string) {
			wg.Go(func() error {
				root := filepath.Join(envy.GoPath(), "src", path, "templates")
				return filepath.Walk(root, func(p string, info os.FileInfo, err error) error {
					if info != nil && !info.IsDir() {
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
							mu.Lock()
							files = append(files, f)
							mu.Unlock()
						}
					}
					return nil
				})
			})
		}(gp)
	}
	err := wg.Wait()
	return files, err
}
