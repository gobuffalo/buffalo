package updater

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/pkg/errors"
)

// MainCheck will either log, or fix, deprecated items in the applications
// main.go file
func MainCheck(*Runner) error {
	fmt.Println("~~~ Checking main.go ~~~")
	b, err := ioutil.ReadFile("main.go")
	if err != nil {
		return errors.WithStack(err)
	}
	if bytes.Contains(b, []byte("app.Start")) {
		fmt.Println("[Warning]: app.Start has been removed in v0.11.0. Use app.Serve Instead.")
	}
	return nil
}
