package resource

import (
	"regexp"

	"github.com/markbates/inflect"
)

// Prop of a model. Starts as name:type on the command line.
type Prop struct {
	Name inflect.Name
	Type string
}

// String representation of Prop
func (m Prop) String() string {
	return string(m.Name)
}

// Valid returns if the property name is valid or not
func (m Prop) Valid() bool {
	reg := regexp.MustCompile(`\A[a-zA-Z]\w+\z`)
	return reg.MatchString(string(m.Name))
}
