package newapp

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gobuffalo/buffalo/meta"
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/pop"
	"github.com/markbates/inflect"
	"github.com/pkg/errors"
)

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
	}
	g.Name = inflect.Name(name)

	if g.Name == "." {
		g.Name = inflect.Name(filepath.Base(g.Root))
	} else {
		g.Root = filepath.Join(g.Root, g.Name.File())
	}

	return g, g.Validate()
}

// Validate that the app generator is good
func (g Generator) Validate() error {
	if g.Name == "" {
		return errors.New("you must enter a name for your new application")
	}

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

	for _, n := range forbiddenAppNames {
		if n == g.Name.Lower() {
			return fmt.Errorf("name %s is not allowed, try a different application name", g.Name)
		}
	}

	if !nameRX.MatchString(string(g.Name)) {
		return fmt.Errorf("name %s is not allowed, application name can only be contain [a-Z0-9-_]", g.Name)
	}

	if s, _ := os.Stat(g.Root); s != nil {
		if !g.Force {
			return fmt.Errorf("%s already exists! Either delete it or use the -f flag to force", g.Name)
		}
	}

	return g.validateInGoPath()
}

func (g Generator) validateInGoPath() error {
	gpMultiple := envy.GoPaths()

	larp := strings.ToLower(g.Root)
	for i := 0; i < len(gpMultiple); i++ {
		lgpm := strings.ToLower(filepath.Join(gpMultiple[i], "src"))
		if strings.HasPrefix(larp, lgpm) {
			return nil
		}
	}

	return ErrNotInGoPath
}

var forbiddenAppNames = []string{"buffalo"}
var nameRX = regexp.MustCompile("^[\\w-]+$")
