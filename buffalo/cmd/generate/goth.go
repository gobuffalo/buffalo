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
			return errors.New("you must specify at least one provider")
		}
		return NewGothGenerator().Run(".", gentronics.Data{
			"providers": args,
		})
	},
}

// NewGothGenerator a actions/goth.go file configured to the specified providers.
func NewGothGenerator() *gentronics.Generator {
	g := gentronics.New()
	g.Add(gentronics.NewFile(filepath.Join("actions", "goth.go"), gGoth))
	g.Add(gentronics.NewCommand(GoGet("github.com/markbates/goth/...")))
	g.Add(Fmt)
	return g
}

var gGoth = `package actions

import (
	"fmt"
	"os"

	"github.com/gobuffalo/buffalo"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	{{#each providers}}
	"github.com/markbates/goth/providers/{{ downcase . }}"
	{{/each}}
)

func init() {
	gothic.Store = App().SessionStore

	goth.UseProviders(
		{{#each providers}}
		{{downcase .}}.New(os.Getenv("{{upcase .}}_KEY"), os.Getenv("{{upcase .}}_SECRET"), fmt.Sprintf("%s%s", App().Host, "/auth/{{downcase .}}/callback")),
		{{/each}}
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
