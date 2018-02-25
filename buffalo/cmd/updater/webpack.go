package updater

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"html/template"

	"github.com/gobuffalo/buffalo/generators/assets/webpack"
	"github.com/gobuffalo/buffalo/generators/newapp"
	"github.com/pkg/errors"
)

// WebpackCheck will compare the current default Buffalo
// webpack.config.js against the applications webpack.config.js. If they are
// different you have the option to overwrite the existing webpack.config.js
// file with the new one.
func WebpackCheck(r *Runner) error {
	fmt.Println("~~~ Checking webpack.config.js ~~~")

	if !r.App.WithWebpack {
		return nil
	}

	g := newapp.Generator{
		App:       r.App,
		Bootstrap: 3,
	}

	box := webpack.TemplateBox

	f, err := box.MustString("webpack.config.js.tmpl")
	if err != nil {
		return errors.WithStack(err)
	}

	tmpl, err := template.New("webpack").Parse(f)
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

	b, err := ioutil.ReadFile("webpack.config.js")
	if err != nil {
		return errors.WithStack(err)
	}

	if string(b) == bb.String() {
		return nil
	}

	if !ask("Your webpack.config.js file is different from the latest Buffalo template.\nWould you like to replace yours with the latest template?") {
		fmt.Println("\tSkipping webpack.config.js")
		return nil
	}

	wf, err := os.Create("webpack.config.js")
	if err != nil {
		return errors.WithStack(err)
	}
	_, err = wf.Write(bb.Bytes())
	if err != nil {
		return errors.WithStack(err)
	}
	return wf.Close()
}
