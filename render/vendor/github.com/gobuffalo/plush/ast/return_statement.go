package ast

import (
	"bytes"

	"github.com/gobuffalo/plush/token"
)

type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

func (rs *ReturnStatement) Printable() bool {
	return true
}

func (rs *ReturnStatement) statementNode() {
}

func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}

func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString("return ")

	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}

	out.WriteString(";")

	return out.String()
}
