package newapp

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gobuffalo/buffalo/meta"
	"github.com/gobuffalo/buffalo/runtime"
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/packr"
	"github.com/gobuffalo/pop"
	"github.com/markbates/inflect"
	"github.com/pkg/errors"
)

// Templates are the templates needed by this generator
var Templates = packr.NewBox("../newapp/templates")

// ErrNotInGoPath can be asserted against
var ErrNotInGoPath = errors.New("currently not in a $GOPATH")

// Generator is the representation of a new Buffalo application
type Generator struct {
	meta.App
	Version     string `json:"version"`
	Force       bool   `json:"force"`
	Verbose     bool   `json:"verbose"`
	DBType      string `json:"db_type"`
	CIProvider  string `json:"ci_provider"`
	AsWeb       bool   `json:"as_web"`
	AsAPI       bool   `json:"as_api"`
	Docker      string `json:"docker"`
	SkipPop     bool   `json:"skip_pop"`
	SkipWebpack bool   `json:"skip_webpack"`
	SkipYarn    bool   `json:"skip_yarn"`
	Bootstrap   int    `json:"bootstrap"`
}

// New app generator
func New(name string) (Generator, error) {
	g := Generator{
		App:        meta.New("."),
		DBType:     "postgres",
		CIProvider: "none",
		AsWeb:      true,
		Docker:     "multi",
		Version:    runtime.Version,
	}
	g.Name = inflect.Name(name)

	if g.Name == "." {
		g.Name = inflect.Name(filepath.Base(g.Root))
	} else {
		g.Root = filepath.Join(g.Root, g.Name.File())
	}

	return g, g.Validate()
}

const header = "Your `buffalo` binary installation is corrupted and is missing vital templates for app creation.\n"
const footer = "Please recheck your installation: https://gobuffalo.io/en/docs/installation."
const goPathAbuse = `It appears you are using multiple GOPATHs:
%s

Using multiple GOPATHs can cause issues with many third party tooling. Please try using only one GOPATH.
`

// ErrTemplatesNotFound means that the `buffalo` binary can't find the template files it needs
// to complete a task. This usually occurs when building `buffalo` locally and/or when using multipath
// GOPATHs.
var ErrTemplatesNotFound = errors.New("templates are missing")

// Validate that the app generator is good
func (g Generator) Validate() error {
	if !Templates.Has("actions/app.go.tmpl") {
		msg := header
		if len(envy.GoPaths()) > 1 {
			msg += "\n" + fmt.Sprintf(goPathAbuse, strings.Join(envy.GoPaths(), "\n"))
		}
		msg += "\n" + footer
		return errors.Wrap(ErrTemplatesNotFound, msg)
	}
	if g.Name == "" {
		return errors.New("you must enter a name for your new application")
	}

	if g.WithPop {
		var found bool
		for _, d := range pop.AvailableDialects {
			if d == g.DBType {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("Unknown db-type %s expecting one of %s", g.DBType, strings.Join(pop.AvailableDialects, ", "))
		}
	}

	for _, n := range forbiddenAppNames {
		if n == g.Name.Lower() {
			return fmt.Errorf("name %s is not allowed, try a different application name", g.Name)
		}
	}

	if !nameRX.MatchString(string(g.Name)) {
		return fmt.Errorf("name %s is not allowed, application name can only contain [a-Z0-9-_]", g.Name)
	}

	if s, _ := os.Stat(g.Root); s != nil {
		if !g.Force {
			return fmt.Errorf("%s already exists! Either delete it or use the -f flag to force", g.Name)
		}
	}

	return g.validateInGoPath()
}

func (g Generator) validateInGoPath() error {
	if g.App.WithModules {
		return nil
	}
	gpMultiple := envy.GoPaths()

	larp := strings.ToLower(meta.ResolveSymlinks(filepath.Dir(g.Root)))

	for i := 0; i < len(gpMultiple); i++ {
		lgpm := strings.ToLower(filepath.Join(gpMultiple[i], "src"))
		if strings.HasPrefix(larp, lgpm) {
			return nil
		}
	}

	return ErrNotInGoPath
}

var forbiddenAppNames = []string{"buffalo"}
var nameRX = regexp.MustCompile(`^[\w-]+$`)
