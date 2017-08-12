package ast

import "github.com/gobuffalo/plush/token"

type HTMLLiteral struct {
	Token token.Token
	Value string
}

func (hl *HTMLLiteral) Printable() bool {
	return true
}

func (hl *HTMLLiteral) expressionNode() {
}

func (hl *HTMLLiteral) TokenLiteral() string {
	return hl.Token.Literal
}

func (hl *HTMLLiteral) String() string {
	return hl.Token.Literal
}
