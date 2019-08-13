package fix

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Plush will update foo.html templates to foo.plush.html templates
func Plush(r *Runner) error {
	templatesDir := filepath.Join(r.App.Root, "templates")
	if _, err := os.Stat(templatesDir); os.IsNotExist(err) {
		// Skip if the templates dir doesn't exist (e.g. API apps)
		return nil
	}
	fmt.Println("~~~ Adding .plush extension to .html/.js/.md files ~~~")
	return filepath.Walk(templatesDir, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		dir := filepath.Dir(p)
		base := filepath.Base(p)

		var exts []string
		ext := filepath.Ext(base)
		for len(ext) != 0 {
			if ext == ".plush" || ext == ".fizz" {
				return nil
			}
			exts = append([]string{ext}, exts...)
			base = strings.TrimSuffix(base, ext)
			ext = filepath.Ext(base)
		}
		exts = append([]string{".plush"}, exts...)

		pn := filepath.Join(dir, base+strings.Join(exts, ""))

		fn, err := os.Create(pn)
		if err != nil {
			return err
		}
		defer fn.Close()

		fo, err := os.Open(p)
		if err != nil {
			return err
		}
		defer fo.Close()
		_, err = io.Copy(fn, fo)

		defer os.Remove(p)

		return err
	})
}
