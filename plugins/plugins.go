package plugins

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/gobuffalo/envy"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// List maps a Buffalo command to a slice of Command
type List map[string]Commands

// Available plugins for the `buffalo` command.
// It will look in $GOPATH/bin and the `./plugins` directory.
// This can be changed by setting the $BUFFALO_PLUGIN_PATH
// environment variable.
//
// Requirements:
// * file/command must be executable
// * file/command must start with `buffalo-`
// * file/command must respond to `available` and return JSON of
//	 plugins.Commands{}
//
// Limit full path scan with direct plugin path
func Available() (List, error) {
	list := List{}
	paths := []string{"plugins"}

	from, err := envy.MustGet("BUFFALO_PLUGIN_PATH")
	if err != nil {
		from, err = envy.MustGet("GOPATH")
		if err != nil {
			return list, errors.WithStack(err)
		}
		from = filepath.Join(from, "bin")
	}

	if runtime.GOOS == "windows" {
		paths = append(paths, strings.Split(from, ";")...)
	} else {
		paths = append(paths, strings.Split(from, ":")...)
	}

	for _, p := range paths {
		if ignorePath(p) {
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
		logrus.Errorf("[PLUGIN] error loading plugin %s: %s\n%s\n", path, err, bb.String())
		return commands
	}
	err = json.NewDecoder(bb).Decode(&commands)
	if err != nil {
		logrus.Errorf("[PLUGIN] error loading plugin %s: %s\n", path, err)
		return commands
	}
	return commands
}

func ignorePath(p string) bool {
	p = strings.ToLower(p)
	for _, x := range []string{`c:\windows`, `c:\program`} {
		if strings.HasPrefix(p, x) {
			return true
		}
	}
	return false
}
