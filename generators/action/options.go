package action

import (
	"errors"

	"github.com/gobuffalo/buffalo/meta"
)

// Options for generating new actions
type Options struct {
	App          meta.App    `json:"app"`
	Name         meta.Name   `json:"name"`
	Method       string      `json:"method"`
	SkipTemplate bool        `json:"skip_template"`
	Actions      []meta.Name `json:"actions"`
	Args         []string    `json:"args"`
}

// NewOptions returns a well formed set of Options
// for generating new actions
func NewOptions(args ...string) (Options, error) {
	o := Options{
		App:     meta.New("."),
		Actions: []meta.Name{},
		Args:    args,
		Method:  "GET",
	}
	if len(args) < 2 {
		return o, errors.New("you need to provide at least an action name and handler name")
	}
	o.Name = meta.Name(args[0])
	for _, a := range args[1:] {
		o.Actions = append(o.Actions, meta.Name(a))
	}

	return o, nil
}
