package ast

import (
	"bytes"
	"strings"

	"github.com/gobuffalo/plush/token"
)

type CallExpression struct {
	Token     token.Token
	Callee    Expression
	Function  Expression
	Arguments []Expression
	Block     *BlockStatement
	ElseBlock *BlockStatement
}

func (ce *CallExpression) expressionNode() {
}

func (ce *CallExpression) TokenLiteral() string {
	return ce.Token.Literal
}

func (ce *CallExpression) String() string {
	var out bytes.Buffer

	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}

	if ce.Callee != nil {
		out.WriteString(ce.Callee.String())
		out.WriteString(".")
	}
	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")
	if ce.Block != nil {
		out.WriteString(" { ")
		out.WriteString(ce.Block.String())
		out.WriteString(" } ")
	}
	if ce.ElseBlock != nil {
		out.WriteString(" else { ")
		out.WriteString(ce.ElseBlock.String())
		out.WriteString(" } ")
	}

	return out.String()
}
