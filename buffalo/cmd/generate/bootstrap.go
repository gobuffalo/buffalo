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
	"path/filepath"

	"github.com/markbates/gentronics"
	"github.com/spf13/cobra"
)

// BootstrapCmd will generate new Bootstrap files. Regardless of whatever
// other settings you might, this will generate jQuery files as that is a
// pre-requisite of Bootstrap.
var BootstrapCmd = &cobra.Command{
	Use:   "bootstrap",
	Short: "Generates Bootstrap 3 files",
	RunE: func(cmd *cobra.Command, args []string) error {
		return NewBootstrapGenerator().Run(".", gentronics.Data{
			"withBootstrap": true,
		})
	},
}

// NewBootstrapGenerator will generate new Bootstrap files. Regardless of whatever
// other settings you might, this will generate jQuery files as that is a
// pre-requisite of Bootstrap.
func NewBootstrapGenerator() *gentronics.Generator {
	should := func(data gentronics.Data) bool {
		if p, ok := data["withBootstrap"]; ok {
			return p.(bool)
		}
		return false
	}
	g := gentronics.New()
	jf := &gentronics.RemoteFile{
		File: gentronics.NewFile(filepath.Join("assets", "css", "bootstrap.css"), ""),
	}
	jf.Should = should
	jf.RemotePath = "https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/css/bootstrap.min.css"
	g.Add(jf)

	jf = &gentronics.RemoteFile{
		File: gentronics.NewFile(filepath.Join("assets", "js", "bootstrap.js"), ""),
	}
	jf.Should = should
	jf.RemotePath = "https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/js/bootstrap.min.js"
	g.Add(jf)
	g.Add(&gentronics.Func{
		Should: should,
		Runner: func(rootPath string, data gentronics.Data) error {
			data["withJQuery"] = true
			return NewJQueryGenerator().Run(rootPath, data)
		},
	})
	return g
}
