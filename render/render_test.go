package render

import (
	"os"
	"path/filepath"

	"github.com/gobuffalo/packr"
)

func withHTMLFile(name string, contents string, fn func(*Engine)) error {
	tmpDir := filepath.Join(os.TempDir(), filepath.Dir(name))
	err := os.MkdirAll(tmpDir, 0766)
	if err != nil {
		return err
	}
	defer os.Remove(tmpDir)

	tmpFile, err := os.Create(filepath.Join(tmpDir, filepath.Base(name)))
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.Write([]byte(contents))
	if err != nil {
		return err
	}

	e := New(Options{
		TemplatesBox: packr.NewBox(os.TempDir()),
	})

	fn(e)
	return nil
}
