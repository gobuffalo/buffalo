package ast

import "github.com/gobuffalo/plush/token"

type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) expressionNode() {
}

func (sl *StringLiteral) TokenLiteral() string {
	return sl.Token.Literal
}

func (sl *StringLiteral) String() string {
	return sl.Token.Literal
}
