package ast

import "github.com/gobuffalo/plush/token"

type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) expressionNode() {
}

func (b *Boolean) TokenLiteral() string {
	return b.Token.Literal
}

func (b *Boolean) String() string {
	return b.Token.Literal
}
