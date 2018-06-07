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

	return filepath.Walk(filepath.Join(r.App.Root, "actions"), func(path string, info os.FileInfo, _ error) error {
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
		// i18n middleware changes in v0.11.1
		if bytes.Contains(b, []byte("T.CookieName")) {
			b = bytes.Replace(b, []byte("T.CookieName"), []byte("T.LanguageExtractorOptions[\"CookieName\"]"), -1)
		}
		if bytes.Contains(b, []byte("T.SessionName")) {
			b = bytes.Replace(b, []byte("T.SessionName"), []byte("T.LanguageExtractorOptions[\"SessionName\"]"), -1)
		}
		if bytes.Contains(b, []byte("T.LanguageFinder=")) || bytes.Contains(b, []byte("T.LanguageFinder ")) {
			r.Warnings = append(r.Warnings, fmt.Sprintf("i18n.Translator#LanguageFinder has been deprecated in v0.11.1. Use i18n.Translator#LanguageExtractors instead. [%s]", path))
		}
		ioutil.WriteFile(path, b, 0664)

		return nil
	})
}
