package translators

import "github.com/markbates/pop/fizz"

type Schema interface {
	TableInfo(string) (*fizz.Table, error)
	Delete(string)
}
