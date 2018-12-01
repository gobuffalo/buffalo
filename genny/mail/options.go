package mail

import (
	"github.com/gobuffalo/meta"
	"github.com/markbates/inflect"
	"github.com/pkg/errors"
)

// Options needed to create a new mailer
type Options struct {
	App      meta.App     `json:"app"`
	Name     inflect.Name `json:"name"`
	SkipInit bool         `json:"skip_init"`
}

// Validate options are useful
func (opts *Options) Validate() error {
	if opts.App.IsZero() {
		opts.App = meta.New(".")
	}

	if len(opts.Name.String()) == 0 {
		return errors.New("you must supply a name for your mailer")
	}
	return nil
}
