package parser

import "github.com/gobuffalo/plush/token"

const (
	_           int = iota
	LOWEST          //
	EQUALS          // ==
	LESSGREATER     // > or <
	SUM             // +
	PRODUCT         // *
	PREFIX          // -X or !X
	CALL            // myFunction(X)
	INDEX           // array[index]
)

var precedences = map[token.Type]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.LTEQ:     LESSGREATER,
	token.GT:       LESSGREATER,
	token.GTEQ:     LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
	token.AND:      PRODUCT,
	token.OR:       PRODUCT,
	token.LPAREN:   CALL,
	token.LBRACKET: INDEX,
}
