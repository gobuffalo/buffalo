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

// RailsJSCmd generates the jQuery UJS file from the Rails project.
var RailsJSCmd = &cobra.Command{
	Use:   "railsjs",
	Short: "Generates an assets/rails.js file",
	Long: `Generates the jQuery UJS file from the Rails project.
More information can be found at:
https://github.com/rails/jquery-ujs`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return NewRailsJSGenerator().Run(".", gentronics.Data{})
	},
}

// NewRailsJSGenerator generates the jQuery UJS file from the Rails project.
func NewRailsJSGenerator() *gentronics.Generator {
	g := gentronics.New()
	jf := &gentronics.RemoteFile{
		File: gentronics.NewFile(filepath.Join("assets", "js", "rails.js"), ""),
	}
	jf.RemotePath = "https://raw.githubusercontent.com/rails/jquery-ujs/master/src/rails.js"
	g.Add(jf)
	return g
}
