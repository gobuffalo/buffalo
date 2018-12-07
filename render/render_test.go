package render

import (
	"github.com/gobuffalo/packr/v2"
	"os"
	"path/filepath"
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
		TemplatesBox: packr.New(os.TempDir(), os.TempDir()),
	})

	fn(e)
	return nil
}
