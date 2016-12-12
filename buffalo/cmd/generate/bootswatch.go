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

package generate

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/markbates/gentronics"
	"github.com/spf13/cobra"
)

// BootswatchCmd will generate new Bootswatch themes. Regardless of whatever
// other settings you have, this will generate jQuery and Bootstrap files as
// they are pre-requisites for Bootswatch
var BootswatchCmd = &cobra.Command{
	Use:   "bootswatch [theme]",
	Short: "Generates Bootswatch 3 files",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("You must choose a theme! [%s]", strings.Join(bootswatchThemes, ", "))
		}
		g, err := NewBootswatchGenerator(args[0])
		if err != nil {
			return err
		}
		return g.Run(".", gentronics.Data{
			"withBootswatch": true,
		})
	},
}

// NewBootswatchGenerator will generate new Bootswatch themes. Regardless of whatever
// other settings you have, this will generate jQuery and Bootstrap files as
// they are pre-requisites for Bootswatch
func NewBootswatchGenerator(theme string) (*gentronics.Generator, error) {
	themeFound := false
	for _, t := range bootswatchThemes {
		if t == theme {
			themeFound = true
			break
		}
	}
	if !themeFound {
		return nil, fmt.Errorf("Could not find a Bootswatch theme for %s!", theme)
	}
	g := gentronics.New()
	g.Add(&gentronics.Func{
		Should: func(data gentronics.Data) bool { return true },
		Runner: func(rootPath string, data gentronics.Data) error {
			data["withJQuery"] = true
			data["withBootstrap"] = true
			err := NewJQueryGenerator().Run(rootPath, data)
			if err != nil {
				return err
			}
			return NewBootstrapGenerator().Run(rootPath, data)
		},
	})
	jf := &gentronics.RemoteFile{
		File: gentronics.NewFile(filepath.Join("assets", "css", "bootstrap.css"), ""),
	}
	jf.RemotePath = fmt.Sprintf("https://maxcdn.bootstrapcdn.com/bootswatch/3.3.7/%s/bootstrap.min.css", theme)
	g.Add(jf)

	return g, nil
}

var bootswatchThemes = []string{"cerulean", "cosmo", "cyborg", "darkly", "flatly", "journal", "lumen", "paper", "readable", "sandstone", "simplex", "slate", "spacelab", "superhero", "united", "yeti"}
