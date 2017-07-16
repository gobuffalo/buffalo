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

	"golang.org/x/sync/errgroup"
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
func Available() (List, error) {
	list := List{}
	paths := []string{"plugins"}
	if runtime.GOOS == "windows" {
		paths = append(paths, strings.Split(os.Getenv("PATH"), ";")...)
	} else {
		paths = append(paths, strings.Split(os.Getenv("PATH"), ":")...)
	}

	ch := make(chan Command)
	wg := &errgroup.Group{}
	for _, p := range paths {
		func(p string) {
			wg.Go(func() error {
				if _, err := os.Stat(p); err != nil {
					return nil
				}
				err := filepath.Walk(p, func(path string, info os.FileInfo, err error) error {
					if info.IsDir() {
						return nil
					}
					base := filepath.Base(path)
					if strings.HasPrefix(base, "buffalo-") {
						wg.Go(func() error {
							return askBin(path, ch)
						})
					}
					return nil
				})
				return err
			})
		}(p)
	}

	go func() {
		for c := range ch {
			bc := c.BuffaloCommand
			if _, ok := list[bc]; !ok {
				list[bc] = Commands{}
			}
			list[bc] = append(list[bc], c)
		}
	}()

	err := wg.Wait()
	close(ch)
	if err != nil {
		return list, errors.WithStack(err)
	}
	return list, nil
}

func askBin(path string, ch chan Command) error {
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
		return nil
	}
	err = json.NewDecoder(bb).Decode(&commands)
	if err != nil {
		fmt.Printf("[PLUGIN] error loading plugin %s: %s\n", path, err)
		return nil
	}
	for _, c := range commands {
		c.Binary = path
		ch <- c
	}
	return nil
}
