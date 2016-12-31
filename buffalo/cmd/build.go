// Copyright © 2016 Mark Bates <mark@markbates.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

var output string

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Builds a Buffalo binary, including bundling of assets (go.rice & webpack)",
	RunE: func(cc *cobra.Command, args []string) error {
		boxes := []string{}
		pwd, err := os.Getwd()
		if err != nil {
			return err
		}
		defer func() {
			for _, b := range boxes {
				fmt.Printf("--> cleaning up rice box %s\n", b)
				os.Remove(b)
			}
		}()

		_, err = os.Stat("webpack.config.js")
		if err == nil {
			// build webpack

			cmd := exec.Command("webpack")
			cmd.Stdin = os.Stdin
			cmd.Stderr = os.Stderr
			cmd.Stdout = os.Stdout
			err = cmd.Run()
			if err != nil {
				return err
			}
		}

		_, err = exec.LookPath("rice")
		if err == nil {
			// if rice exists, try and build some boxes:
			err = filepath.Walk(pwd, func(path string, info os.FileInfo, err error) error {
				if info.IsDir() {
					if filepath.Base(path) == "node_modules" {
						return filepath.SkipDir
					}
					err = os.Chdir(path)
					if err != nil {
						return err
					}
					cmd := exec.Command("rice", "embed-go")
					err = cmd.Run()
					if err == nil {
						bp := filepath.Join(path, "rice-box.go")
						_, err := os.Stat(bp)
						if err == nil {
							fmt.Printf("--> built rice box %s\n", bp)
							boxes = append(boxes, bp)
						}
					}
				}
				return nil
			})
			if err != nil {
				return err
			}
		}

		os.Chdir(pwd)
		cmd := exec.Command("go", "build", "-v", "-o", output)
		cmd.Stdin = os.Stdin
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		return cmd.Run()

	},
}

func init() {
	RootCmd.AddCommand(buildCmd)
	pwd, _ := os.Getwd()
	buildCmd.Flags().StringVarP(&output, "output", "o", filepath.Base(pwd), "set the name of the binary")
}
