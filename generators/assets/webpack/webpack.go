package webpack

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gobuffalo/buffalo/generators/assets"
	"github.com/gobuffalo/buffalo/generators/common"
	"github.com/markbates/gentronics"
)

var logo = &gentronics.RemoteFile{
	File:       gentronics.NewFile("assets/images/logo.svg", ""),
	RemotePath: assets.LogoURL,
}

// BinPath is the path to the local install of webpack
var BinPath = filepath.Join("node_modules", ".bin", "webpack")

// New webpack generator
func New(data gentronics.Data) (*gentronics.Generator, error) {
	g := gentronics.New()

	// if there's no npm, return!
	_, err := exec.LookPath("npm")
	if err != nil {
		fmt.Println("Could not find npm/node. Skipping webpack generation.")
		return g, nil
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

	files, err := common.Find(filepath.Join("assets", "webpack"))
	if err != nil {
		return g, err
	}

	for _, f := range files {
		g.Add(gentronics.NewFile(f.WritePath, f.Body))
	}

	c := gentronics.NewCommand(exec.Command(command, "init", "-y"))
	g.Add(c)

	modules := []string{"webpack@^2.2.1", "sass-loader", "css-loader", "style-loader", "node-sass",
		"babel-loader", "extract-text-webpack-plugin", "babel", "babel-core", "url-loader", "file-loader",
		"jquery", "bootstrap", "path", "font-awesome", "npm-install-webpack-plugin", "jquery-ujs",
		"copy-webpack-plugin", "expose-loader",
	}

	args = append(args, modules...)
	g.Add(gentronics.NewCommand(exec.Command(command, args...)))
	return g, nil
}

func generateYarn(data gentronics.Data) error {
	// if there's no yarn, install it!
	_, err := exec.LookPath("yarn")
	// A new gentronics is necessary to have yarn available in path
	if err != nil {
		yg := gentronics.New()
		yargs := []string{"install", "-g", "yarn"}
		yg.Add(gentronics.NewCommand(exec.Command("npm", yargs...)))
		err = yg.Run(".", data)
		if err != nil {
			return err
		}
	}
	return nil
}
