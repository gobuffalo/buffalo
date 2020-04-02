package ci

import (
	"fmt"
	"strings"

	"github.com/gobuffalo/buffalo/runtime"

	"github.com/gobuffalo/meta"
	"github.com/gobuffalo/pop/v5"
)

// Available CI implementations
var Available = []string{"travis", "gitlab"}

// Options for CI
type Options struct {
	App      meta.App
	DBType   string
	Provider string
	Version  string
}

// Validate that options are usuable
func (opts *Options) Validate() error {
	if opts.App.IsZero() {
		opts.App = meta.New(".")
	}

	if len(opts.Version) == 0 {
		opts.Version = runtime.Version
	}

	if len(opts.Provider) == 0 {
		return fmt.Errorf("no provider chosen")
	}
	opts.Provider = strings.ToLower(opts.Provider)

	var found bool
	for _, a := range Available {
		if opts.Provider == a {
			found = true
			break
		}
		if opts.Provider == a+"-ci" {
			opts.Provider = a
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("unknown provider %s expecting one of %s", opts.Provider, strings.Join(Available, ", "))
	}

	found = false
	for _, d := range pop.AvailableDialects {
		if d == opts.DBType {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("unknown dialect %q expecting one of %s", opts.DBType, strings.Join(pop.AvailableDialects, ", "))
	}
	return nil
}
