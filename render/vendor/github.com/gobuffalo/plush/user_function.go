package plush

import (
	"bytes"
	"strings"

	"github.com/gobuffalo/plush/ast"
)

type userFunction struct {
	Parameters []*ast.Identifier
	Block      *ast.BlockStatement
}

func (f *userFunction) String() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Block.String())
	out.WriteString("\n}")

	return out.String()
}
