package resource

import (
	"github.com/gobuffalo/attrs"
	"github.com/gobuffalo/flect/name"
	"github.com/gobuffalo/meta"
)

type presenter struct {
	App   meta.App
	Name  name.Ident
	Model name.Ident
	Attrs attrs.Attrs
}
