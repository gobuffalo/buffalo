package actions

import "github.com/gobuffalo/flect/name"

type data map[string]interface{}

type presenter struct {
	Name    name.Ident
	Actions []name.Ident
	Helpers data
	Data    data
	Options *Options
}
