package resource

import (
	"strings"

	"github.com/gobuffalo/buffalo/meta"
	"github.com/markbates/inflect"
)

// Prop of a model. Starts as name:type on the command line.
type Prop struct {
	Name meta.Name
	Type string
}

// String representation of Prop
func (m Prop) String() string {
	return string(m.Name)
}

func modelPropertiesFromArgs(args []string) []Prop {
	var props []Prop
	if len(args) == 0 {
		return props
	}
	for _, a := range args[1:] {
		ax := strings.Split(a, ":")
		p := Prop{
			Name: meta.Name(inflect.ForeignKeyToAttribute(ax[0])),
			Type: "string",
		}
		if len(ax) > 1 {
			p.Type = strings.ToLower(strings.TrimPrefix(ax[1], "nulls."))
		}
		props = append(props, p)
	}
	return props
}
