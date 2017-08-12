package ast

import (
	"bytes"
	"github.com/gobuffalo/plush/token"
	"strings"
)

type FunctionLiteral struct {
	Token      token.Token
	Parameters []*Identifier
	Block      *BlockStatement
}

func (fl *FunctionLiteral) expressionNode() {
}

func (fl *FunctionLiteral) TokenLiteral() string {
	return fl.Token.Literal
}

func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}

	out.WriteString(fl.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(fl.Block.String())

	return out.String()
}
