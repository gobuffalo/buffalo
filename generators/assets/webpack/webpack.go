package webpack

import (
	"os/exec"
	"path/filepath"

	"github.com/gobuffalo/buffalo/generators"
	"github.com/gobuffalo/buffalo/generators/assets"
	"github.com/gobuffalo/buffalo/generators/assets/standard"
	"github.com/gobuffalo/makr"
	"github.com/gobuffalo/packr"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// TemplateBox contains all templates needed for the webpack generator
var TemplateBox = packr.NewBox("../webpack/templates")

var logo = &makr.RemoteFile{
	File:       makr.NewFile("assets/images/logo.svg", ""),
	RemotePath: assets.LogoURL,
}

// BinPath is the path to the local install of webpack
var BinPath = filepath.Join("node_modules", ".bin", "webpack")

// Run webpack generator
func (w Generator) Run(root string, data makr.Data) error {
	g := makr.New()

	// if there's no npm, return!
	if _, err := exec.LookPath("npm"); err != nil {
		logrus.Info("Could not find npm. Skipping webpack generation.")

		return standard.Run(root, data)
	}

	command := "yarn"

	if !w.WithYarn {
		command = "npm"
	} else {
		err := installYarn(data)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	g.Add(logo)

	files, err := generators.FindByBox(TemplateBox)
	if err != nil {
		return errors.WithStack(err)
	}

	for _, f := range files {
		g.Add(makr.NewFile(f.WritePath, f.Body))
	}

	args := []string{"install", "--no-progress", "--save"}
	g.Add(makr.NewCommand(exec.Command(command, args...)))
	data["opts"] = w
	return g.Run(root, data)
}

func installYarn(data makr.Data) error {
	// if there's no yarn, install it!
	_, err := exec.LookPath("yarn")
	// A new makr is necessary to have yarn available in path
	if err != nil {
		yg := makr.New()
		yargs := []string{"install", "-g", "yarn"}
		yg.Add(makr.NewCommand(exec.Command("npm", yargs...)))
		err = yg.Run(".", data)
		if err != nil {
			return err
		}
	}
	return nil
}
