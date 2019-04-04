package buffalo

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"text/tabwriter"

	"github.com/gobuffalo/x/randx"
	"github.com/markbates/grift/grift"
	"golang.org/x/crypto/bcrypt"
)

// Grifts decorates the app with tasks
func Grifts(app *App) {
	routesGrift(app)
	middlewareGrift(app)
	secretGrift()
}

func secretGrift() {
	grift.Desc("secret", "Generate a cryptographically secure secret key")
	grift.Add("secret", func(c *grift.Context) error {
		bb := []byte{}
		for i := 0; i < 4; i++ {
			b := []byte(randx.String(64))
			b, err := bcrypt.GenerateFromPassword(b, bcrypt.DefaultCost)
			if err != nil {
				return err
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
	printMiddlewareByRoute(a)
}

func printMiddlewareByRoute(a *App) {
	mws := map[string]string{}
	for _, r := range a.Routes() {
		if mws[r.App.Name] == "" {
			pname := ""
			if parent := getParentApp(r.App.root, r.App.Name); parent != nil {
				pname = parent.Name
			}

			mws[r.App.Name] = r.App.Middleware.String()
			if mws[pname] != mws[r.App.Name] {
				fmt.Printf("-> %s\n", r.App.Name)
				printMiddlewareStackWithIndent(mws[r.App.Name])
			} else {
				fmt.Printf("-> %s (see: %v)\n", r.App.Name, pname)
			}
		}
		s := "\n" + mws[r.App.Name]
		for k := range r.App.Middleware.skips {
			mw := strings.Split(k, funcKeyDelimeter)[0]
			h := strings.Split(k, funcKeyDelimeter)[1]
			if h == r.HandlerName {
				s = strings.Replace(s, "\n"+mw, "", 1)
			}
		}
		if "\n"+mws[r.App.Name] != s {
			ahn := strings.Split(r.HandlerName, "/")
			hn := ahn[len(ahn)-1]
			fmt.Printf("-> %s %s (by %s)\n", r.Method, r.Path, hn)
			printMiddlewareStackWithIndent(s)
		}
	}
}

func getParentApp(r *App, name string) *App {
	if r == nil {
		return nil
	}
	for _, x := range r.children {
		if x.Name == name {
			return r
		}
		if len(x.children) > 0 {
			if ret := getParentApp(x, name); ret != nil {
				return ret
			}
		}
	}
	return nil
}

func printMiddlewareStackWithIndent(s string) {
	if s == "" {
		s = "[none]"
	}
	s = strings.Replace(s, "\n", "\n\t", -1)
	fmt.Printf("\t%v\n", strings.TrimSpace(s))
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
