package action

import (
	"errors"

	"github.com/gobuffalo/buffalo/meta"
)

// Generator for generating new actions
type Generator struct {
	App          meta.App    `json:"app"`
	Name         meta.Name   `json:"name"`
	Method       string      `json:"method"`
	SkipTemplate bool        `json:"skip_template"`
	Actions      []meta.Name `json:"actions"`
	Args         []string    `json:"args"`
}

// New returns a well formed set of Options
// for generating new actions
func New(args ...string) (Generator, error) {
	o := Generator{
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
