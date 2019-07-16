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
	fmt.Println("~~~ Adding .plush extension to .html files ~~~")
	return filepath.Walk(r.App.Root, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		dir := filepath.Dir(p)
		base := filepath.Base(p)
		sep := strings.Split(base, ".")

		if filepath.Ext(p) != ".html" {
			return nil
		}

		if len(sep) != 2 {
			return nil
		}

		pn := filepath.Join(dir, sep[0]+".plush."+sep[1])

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
