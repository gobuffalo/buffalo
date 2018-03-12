package updater

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

// DeprecrationsCheck will either log, or fix, deprecated items in the application
func DeprecrationsCheck(r *Runner) error {
	fmt.Println("~~~ Checking for deprecations ~~~")
	b, err := ioutil.ReadFile("main.go")
	if err != nil {
		return errors.WithStack(err)
	}
	if bytes.Contains(b, []byte("app.Start")) {
		r.Warnings = append(r.Warnings, "app.Start has been removed in v0.11.0. Use app.Serve Instead. [main.go]")
	}

	return filepath.Walk(filepath.Join(r.App.Root, "actions"), func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if filepath.Ext(path) != ".go" {
			return nil
		}

		b, err := ioutil.ReadFile(path)
		if err != nil {
			return errors.WithStack(err)
		}
		if bytes.Contains(b, []byte("Websocket()")) {
			r.Warnings = append(r.Warnings, fmt.Sprintf("buffalo.Context#Websocket has been deprecated in v0.11.0. Use github.com/gorilla/websocket directly. [%s]", path))
		}

		return nil
	})
}
