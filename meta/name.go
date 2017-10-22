package meta

import (
	"strings"

	"github.com/markbates/inflect"
)

// Name is a string that represents the "name" of a thing, like an app, model, etc...
type Name string

// Title version of a name. ie. "foo_bar" => "Foo Bar"
func (n Name) Title() string {
	x := strings.Split(string(n), "/")
	for i, s := range x {
		x[i] = inflect.Titleize(s)
	}

	return strings.Join(x, " ")
}

// Underscore version of a name. ie. "FooBar" => "foo_bar"
func (n Name) Underscore() string {
	return inflect.Underscore(string(n))
}

// Plural version of a name
func (n Name) Plural() string {
	return inflect.Pluralize(string(n))
}

// Singular version of a name
func (n Name) Singular() string {
	return inflect.Singularize(string(n))
}

// Camel version of a name
func (n Name) Camel() string {
	return inflect.Camelize(string(n))
}

// Model version of a name. ie. "user" => "User"
func (n Name) Model() string {
	x := strings.Split(string(n), "/")
	for i, s := range x {
		x[i] = inflect.Singularize(inflect.Camelize(s))
	}

	return strings.Join(x, "")
}

// Resource version of a name
func (n Name) Resource() string {
	x := strings.Split(string(n), "/")
	for i, s := range x {
		x[i] = inflect.Camelize(s)
	}

	return inflect.Pluralize(strings.Join(x, ""))
}

// ModelPlural version of a name. ie. "user" => "Users"
func (n Name) ModelPlural() string {
	return inflect.Pluralize(n.Model())
}

// File version of a name
func (n Name) File() string {
	return inflect.Underscore(inflect.Camelize(string(n)))
}

// Table version of a name
func (n Name) Table() string {
	return inflect.Underscore(inflect.Pluralize(string(n)))
}

// UnderSingular version of a name
func (n Name) UnderSingular() string {
	return inflect.Underscore(inflect.Singularize(string(n)))
}

// PluralCamel version of a name
func (n Name) PluralCamel() string {
	return inflect.Pluralize(inflect.Camelize(string(n)))
}

// PluralUnder version of a name
func (n Name) PluralUnder() string {
	return inflect.Pluralize(inflect.Underscore(string(n)))
}

// URL version of a name
func (n Name) URL() string {
	return n.PluralUnder()
}

// CamelSingular version of a name
func (n Name) CamelSingular() string {
	return inflect.Camelize(inflect.Singularize(string(n)))
}

// VarCaseSingular version of a name. ie. "FooBar" => "fooBar"
func (n Name) VarCaseSingular() string {
	return inflect.CamelizeDownFirst(inflect.Singularize(n.Resource()))
}

// VarCasePlural version of a name. ie. "FooBar" => "fooBar"
func (n Name) VarCasePlural() string {
	return inflect.CamelizeDownFirst(inflect.Pluralize(n.Resource()))
}

// Lower case version of a string
func (n Name) Lower() string {
	return strings.ToLower(string(n))
}
