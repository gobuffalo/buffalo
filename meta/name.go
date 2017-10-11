package meta

import "github.com/markbates/inflect"

// Name is a string that represents the "name" of a thing, like an app, model, etc...
type Name string

// Title version of a name. ie. "foo_bar" => "Foo Bar"
func (n Name) Title() string {
	return inflect.Titleize(string(n))
}

// Underscore version of a name. ie. "FooBar" => "foo_bar"
func (n Name) Underscore() string {
	return inflect.Underscore(string(n))
}
