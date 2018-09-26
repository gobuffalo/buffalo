package buffalo

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"text/tabwriter"

	"github.com/gobuffalo/buffalo/buffalo/cmd"
	"github.com/gobuffalo/x/randx"
	"github.com/markbates/grift/grift"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// Grifts decorates the app with tasks
func Grifts(app *App) {
	routesGrift(app)
	middlewareGrift(app)
	secretGrift()
	versionGrift()
}

// version this is here EXPLICITLY to force
// the cmd and runtime packages to be imported
// as they're required for building and tools
// like dep prune those packages because they're
// not imported. DO NOT REMOVE!!! :)
// https://github.com/gobuffalo/buffalo/pull/1325
func versionGrift() {
	grift.Namespace("buffalo", func() {
		grift.Desc("version", "Print the version number of buffalo")
		grift.Add("version", func(c *grift.Context) error {
			cmd.VersionCmd.Run(nil, []string{})
			return nil
		})
	})
}

func secretGrift() {
	grift.Desc("secret", "Generate a cryptographically secure secret key")
	grift.Add("secret", func(c *grift.Context) error {
		bb := []byte{}
		for i := 0; i < 4; i++ {
			b := []byte(randx.String(64))
			b, err := bcrypt.GenerateFromPassword(b, bcrypt.DefaultCost)
			if err != nil {
				return errors.WithStack(err)
			}
			bb = append(bb, b...)
		}
		rx := regexp.MustCompile(`(\W+)`)
		bb = rx.ReplaceAll(bb, []byte(""))
		s := randx.String(6) + string(bb)
		fmt.Println(s[:127])
		return nil
	})
}

func middlewareGrift(a *App) {
	grift.Desc("middleware", "Prints out your middleware stack")
	grift.Add("middleware", func(c *grift.Context) error {
		printMiddleware(a)
		return nil
	})
}

func printMiddleware(a *App) {
	fmt.Printf("-> %s\n", a.Name)
	fmt.Printf("%v\n", a.Middleware.String())
	for _, x := range a.children {
		printMiddleware(x)
	}
}

func routesGrift(a *App) {
	grift.Desc("routes", "Print out all defined routes")
	grift.Add("routes", func(c *grift.Context) error {
		routes := a.Routes()
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)
		fmt.Fprintln(w, "METHOD\t PATH\t ALIASES\t NAME\t HANDLER")
		fmt.Fprintln(w, "------\t ----\t -------\t ----\t -------")
		for _, r := range routes {
			fmt.Fprintf(w, "%s\t %s\t %s\t %s\t %s\n", r.Method, r.Path, strings.Join(r.Aliases, " "), r.PathName, r.HandlerName)
		}
		w.Flush()
		return nil
	})
}
