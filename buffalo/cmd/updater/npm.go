package updater

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"html/template"

	"github.com/gobuffalo/buffalo/generators/assets/webpack"
	"github.com/gobuffalo/buffalo/generators/newapp"
	"github.com/pkg/errors"
)

// PackageJSONCheck will compare the current default Buffalo
// package.json against the applications package.json. If they are
// different you have the option to overwrite the existing package.json
// file with the new one.
func PackageJSONCheck(r *Runner) error {
	fmt.Println("~~~ Checking package.json ~~~")

	if !r.App.WithWebpack {
		return nil
	}

	g := newapp.Generator{
		App:       r.App,
		Bootstrap: 3,
	}

	box := webpack.TemplateBox

	f, err := box.MustString("package.json.tmpl")
	if err != nil {
		return errors.WithStack(err)
	}

	tmpl, err := template.New("package.json").Parse(f)
	if err != nil {
		return errors.WithStack(err)
	}

	bb := &bytes.Buffer{}
	err = tmpl.Execute(bb, map[string]interface{}{
		"opts": g,
	})
	if err != nil {
		return errors.WithStack(err)
	}

	b, err := ioutil.ReadFile("package.json")
	if err != nil {
		return errors.WithStack(err)
	}

	if string(b) == bb.String() {
		return nil
	}

	if !ask("Your package.json file is different from the latest Buffalo template.\nWould you like to REPLACE yours with the latest template?") {
		fmt.Println("\tskipping package.json")
		return nil
	}

	pf, err := os.Create("package.json")
	if err != nil {
		return errors.WithStack(err)
	}
	_, err = pf.Write(bb.Bytes())
	if err != nil {
		return errors.WithStack(err)
	}
	err = pf.Close()
	if err != nil {
		return errors.WithStack(err)
	}

	var cmd *exec.Cmd
	if r.App.WithYarn {
		cmd = exec.Command("yarn", "install")
	} else {
		cmd = exec.Command("npm", "install")
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
