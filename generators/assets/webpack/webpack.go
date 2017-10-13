package webpack

import (
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/gobuffalo/buffalo/generators"
	"github.com/gobuffalo/buffalo/generators/assets"
	"github.com/gobuffalo/buffalo/generators/assets/standard"
	"github.com/gobuffalo/makr"
	"github.com/pkg/errors"
)

var logo = &makr.RemoteFile{
	File:       makr.NewFile("assets/images/logo.svg", ""),
	RemotePath: assets.LogoURL,
}

// BinPath is the path to the local install of webpack
var BinPath = filepath.Join("node_modules", ".bin", "webpack")

// Run webpack generator
func Run(root string, data makr.Data) error {
	g := makr.New()

	command := "yarn"

	if b, ok := data["withYarn"].(bool); ok && !b {
		command = "npm"
	} else {
		err := installYarn(data)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	// if there's no npm, return!
	_, err := exec.LookPath("npm")
	if err != nil {
		fmt.Println("Could not find npm. Skipping webpack generation.")

		return standard.Run(root, data)
	}

	g.Add(logo)

	files, err := generators.Find(filepath.Join(generators.TemplatesPath, "assets", "webpack"))
	if err != nil {
		return errors.WithStack(err)
	}

	for _, f := range files {
		g.Add(makr.NewFile(f.WritePath, f.Body))
	}

	args := []string{"install", "--no-progress", "--save"}
	g.Add(makr.NewCommand(exec.Command(command, args...)))
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
