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

// JQueryCmd will generate jQuery files.
var JQueryCmd = &cobra.Command{
	Use:   "jquery",
	Short: "Generates an assets/jquery.js file",
	RunE: func(cmd *cobra.Command, args []string) error {
		data := gentronics.Data{
			"withJQuery": true,
		}
		return NewJQueryGenerator().Run(".", data)
	},
}

// NewJQueryGenerator will generate jQuery files.
func NewJQueryGenerator() *gentronics.Generator {
	should := func(data gentronics.Data) bool {
		if p, ok := data["withJQuery"]; ok {
			return p.(bool)
		}
		return false
	}

	g := gentronics.New()
	jf := &gentronics.RemoteFile{
		File: gentronics.NewFile(filepath.Join("assets", "js", "jquery.js"), ""),
	}
	jf.Should = should
	jf.RemotePath = "https://cdnjs.cloudflare.com/ajax/libs/jquery/3.1.1/jquery.min.js"
	g.Add(jf)

	jm := &gentronics.RemoteFile{
		File: gentronics.NewFile(filepath.Join("assets", "js", "jquery.map"), ""),
	}
	jm.Should = should
	jm.RemotePath = "https://cdnjs.cloudflare.com/ajax/libs/jquery/3.1.1/jquery.min.map"
	g.Add(jm)
	return g
}
