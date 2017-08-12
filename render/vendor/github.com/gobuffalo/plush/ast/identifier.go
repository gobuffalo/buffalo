package ast

import (
	"bytes"

	"github.com/gobuffalo/plush/token"
)

type Identifier struct {
	Token  token.Token
	Callee *Identifier
	Value  string
}

func (i *Identifier) expressionNode() {
}

func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

func (i *Identifier) String() string {
	out := &bytes.Buffer{}
	if i.Callee != nil {
		out.WriteString(i.Callee.String())
		out.WriteString(".")
	}
	out.WriteString(i.Value)
	return out.String()
}
