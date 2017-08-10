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

// New webpack generator
func New(data makr.Data) (*makr.Generator, error) {
	g := makr.New()

	// if there's no npm, return!
	_, err := exec.LookPath("npm")
	if err != nil {
		fmt.Println("Could not find npm/node. Skipping webpack generation.")

		wg, err := standard.New(data)
		if err != nil {
			return g, errors.WithStack(err)
		}
		return wg, nil
	}

	command := "npm"
	args := []string{"install", "--no-progress", "--save"}
	if b, ok := data["withYarn"].(bool); ok && b {
		command = "yarn"
		args = []string{"add"}
		err := generateYarn(data)
		if err != nil {
			return g, err
		}
	}

	g.Add(logo)

	files, err := generators.Find(filepath.Join("assets", "webpack"))
	if err != nil {
		return g, err
	}

	for _, f := range files {
		g.Add(makr.NewFile(f.WritePath, f.Body))
	}

	c := makr.NewCommand(exec.Command(command, "init", "--no-progress", "-y"))
	g.Add(c)

	modules := []string{
		"webpack@~2.3.0",
		"sass-loader@~6.0.5",
		"css-loader@~0.28.4",
		"expose-loader@~0.7.3",
		"style-loader@~0.18.2",
		"node-sass@~4.5.3",
		"extract-text-webpack-plugin@2.1.2",
		"babel-cli@~6.24.1",
		"babel-core@~6.25.0",
		"babel-preset-env@~1.5.2",
		"babel-loader@~7.0.0",
		"url-loader@~0.5.9",
		"file-loader@~0.11.2",
		"jquery@~3.2.1",
		"bootstrap@~3.3.7",
		"path@~0.12.7",
		"font-awesome@~4.7.0",
		"npm-install-webpack-plugin@4.0.4",
		"jquery-ujs@~1.2.2",
		"copy-webpack-plugin@~4.0.1",
		"uglifyjs-webpack-plugin@~0.4.6",
		"webpack-manifest-plugin@~1.2.1",
	}

	args = append(args, modules...)
	g.Add(makr.NewCommand(exec.Command(command, args...)))
	return g, nil
}

func generateYarn(data makr.Data) error {
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
