package meta

import "github.com/markbates/inflect"

type Name string

func (n Name) Title() string {
	return inflect.Titleize(string(n))
}

func (n Name) Underscore() string {
	return inflect.Underscore(string(n))
}
