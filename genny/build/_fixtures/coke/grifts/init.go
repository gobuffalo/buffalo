package grifts

import (
	"github.com/gobuffalo/buffalo"
	"github.com/markbates/coke/actions"
)

func init() {
	buffalo.Grifts(actions.App())
}
