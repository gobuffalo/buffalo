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
	"errors"
	"path/filepath"

	"github.com/markbates/gentronics"
	"github.com/spf13/cobra"
)

// GothCmd generates a actions/goth.go file configured to the specified providers.
var GothCmd = &cobra.Command{
	Use:   "goth [provider provider...]",
	Short: "Generates a actions/goth.go file configured to the specified providers.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("You must specifiy at least one provider!")
		}
		return NewGothGenerator().Run(".", gentronics.Data{
			"providers": args,
		})
	},
}

// NewGothGenerator a actions/goth.go file configured to the specified providers.
func NewGothGenerator() *gentronics.Generator {
	g := gentronics.New()
	f := gentronics.NewFile(filepath.Join("actions", "goth.go"), gGoth)
	g.Add(f)
	g.Add(Fmt)
	return g
}

var gGoth = `package actions

import (
	"fmt"
	"os"

	"github.com/markbates/buffalo"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	{{ range .providers -}}
	"github.com/markbates/goth/providers/{{. | downcase}}"
	{{ end -}}
)

func init() {
	gothic.Store = App().SessionStore

	goth.UseProviders(
		{{ range .providers -}}
		{{.|downcase}}.New(os.Getenv("{{.|upcase}}_KEY"), os.Getenv("{{.|upcase}}_SECRET"), fmt.Sprintf("%s%s", App().Host, "/auth/{{.|downcase}}/callback")),
		{{ end -}}
	)

	app := App().Group("/auth")
	app.GET("/{provider}", buffalo.WrapHandlerFunc(gothic.BeginAuthHandler))
	app.GET("/{provider}/callback", AuthCallback)
}

func AuthCallback(c buffalo.Context) error {
	user, err := gothic.CompleteUserAuth(c.Response(), c.Request())
	if err != nil {
		return c.Error(401, err)
	}
	// Do something with the user, maybe register them/sign them in
	return c.Render(200, r.JSON(user))
}
`
