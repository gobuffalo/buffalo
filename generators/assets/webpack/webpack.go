package webpack

import (
	"fmt"
	"os"
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
	args := []string{"install", "--save"}
	// If yarn.lock exists then yarn is used by default (generate webpack)
	_, ferr := os.Stat("yarn.lock")
	if ferr == nil {
		data["withYarn"] = true
	}

	if _, ok := data["withYarn"]; ok {
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

	c := makr.NewCommand(exec.Command(command, "init", "-y"))
	g.Add(c)

	modules := []string{"webpack@~2.3.0", "sass-loader", "css-loader", "style-loader", "node-sass",
		"extract-text-webpack-plugin@2.1.2", "babel-cli", "babel-core", "babel-preset-env", "babel-loader", "url-loader",
		"file-loader", "jquery", "bootstrap", "path", "font-awesome", "npm-install-webpack-plugin", "jquery-ujs",
		"copy-webpack-plugin", "expose-loader", "uglifyjs-webpack-plugin@~0.4.6",
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
