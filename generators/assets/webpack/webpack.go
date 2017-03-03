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

// BinPath is the path to the local install of webpack
var BinPath = filepath.Join("node_modules", ".bin", "webpack")

func New(data gentronics.Data) (*gentronics.Generator, error) {
	g := gentronics.New()

	should := func(data gentronics.Data) bool {
		return true
	}

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

	useYarn := func(data gentronics.Data) bool {
		if b, ok := data["withYarn"]; ok {
			return b.(bool)
		}
		return false
	}
	if useYarn(data) {
		// if there's no yarn, install it!
		_, err := exec.LookPath("yarn")
		// A new gentronics is necessary to have yarn available in path
		if err != nil {
			yg := gentronics.New()
			yargs := []string{"install", "-g", "yarn"}
			yg.Should = useYarn
			yg.Add(gentronics.NewCommand(exec.Command(command, yargs...)))
			err = yg.Run(".", data)
			if err != nil {
				return g, err
			}
		}
		command = "yarn"
		args = []string{"add"}
	}

	g.Should = should
	g.Add(assets.AssetsLogo)

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
