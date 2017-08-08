package plugins

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// List maps a Buffalo command to a slice of Command
type List map[string]Commands

// Available plugins for the `buffalo` command.
// It will look in $PATH and the `./plugins` directory.
//
// Requirements:
// * file/command must be executable
// * file/command must start with `buffalo-`
// * file/command must respond to `available` and return JSON of
//	 plugins.Commands{}
//
// Caveats:
// * The C:\Windows directory is excluded
func Available() (List, error) {
	list := List{}
	paths := []string{"plugins"}
	if runtime.GOOS == "windows" {
		paths = append(paths, strings.Split(os.Getenv("PATH"), ";")...)
	} else {
		paths = append(paths, strings.Split(os.Getenv("PATH"), ":")...)
	}
	for _, p := range paths {
		if strings.HasPrefix(strings.ToLower(p), `c:\windows`) {
			continue	
		}
		if _, err := os.Stat(p); err != nil {
			continue
		}
		err := filepath.Walk(p, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				// May indicate a permissions problem with the path, skip it
				return nil
			}
			if info.IsDir() {
				return nil
			}
			base := filepath.Base(path)
			if strings.HasPrefix(base, "buffalo-") {
				commands := askBin(path)
				for _, c := range commands {
					bc := c.BuffaloCommand
					if _, ok := list[bc]; !ok {
						list[bc] = Commands{}
					}
					c.Binary = path
					list[bc] = append(list[bc], c)
				}
			}
			return nil
		})
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}
	return list, nil
}

func askBin(path string) Commands {
	commands := Commands{}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, path, "available")
	bb := &bytes.Buffer{}
	cmd.Stdout = bb
	cmd.Stderr = bb
	err := cmd.Run()
	if err != nil {
		fmt.Printf("[PLUGIN] error loading plugin %s: %s\n%s\n", path, err, bb.String())
		return commands
	}
	err = json.NewDecoder(bb).Decode(&commands)
	if err != nil {
		fmt.Printf("[PLUGIN] error loading plugin %s: %s\n", path, err)
		return commands
	}
	return commands
}
