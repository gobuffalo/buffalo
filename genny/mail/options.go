package mail

import (
	"github.com/gobuffalo/buffalo/meta"
	"github.com/markbates/inflect"
)

// Options needed to create a new mailer
type Options struct {
	App      meta.App     `json:"app"`
	Name     inflect.Name `json:"name"`
	SkipInit bool         `json:"skip_init"`
}
