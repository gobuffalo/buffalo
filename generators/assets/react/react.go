package react

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gobuffalo/buffalo/generators"
	"github.com/gobuffalo/buffalo/generators/assets"
	"github.com/gobuffalo/makr"
)

var logo = &makr.RemoteFile{
	File:       makr.NewFile("assets/logo.svg", ""),
	RemotePath: assets.LogoURL,
}

// New react generator
func New(data makr.Data) (*makr.Generator, error) {
	g := makr.New()

	// if there's no npm, return!
	_, err := exec.LookPath("npm")
	if err != nil {
		fmt.Println("Could not find npm/node. Skipping react generation.")
		return g, nil
	}

	command := "npm"
	args := []string{"install"}
	// If yarn.lock exists then yarn is used by default (generate react)
	_, ferr := os.Stat("yarn.lock")
	if ferr == nil {
		data["withYarn"] = true
	}

	if _, ok := data["withYarn"]; ok {
		command = "yarn"
		err := generateYarn(data)
		if err != nil {
			return g, err
		}
	}

	g.Add(logo)

	files, err := generators.Find(filepath.Join("assets", "react"))
	if err != nil {
		return g, err
	}

	for _, f := range files {
		g.Add(makr.NewFile(f.WritePath, f.Body))
	}

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
