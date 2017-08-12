package ast

import "github.com/gobuffalo/plush/token"

type IntegerLiteral struct {
	Token token.Token
	Value int
}

func (il *IntegerLiteral) expressionNode() {
}

func (il *IntegerLiteral) TokenLiteral() string {
	return il.Token.Literal
}

func (il *IntegerLiteral) String() string {
	return il.Token.Literal
}
