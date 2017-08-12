package token

// Type represents each type of token.
type Type string

// Token of a section of input source.
type Token struct {
	Type    Type
	Literal string
}

var keywords = map[string]Type{
	"fn":     FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
	"for":    FOR,
	"in":     IN,
}

// LookupIdent an ident and return a keyword type, or a plain ident
func LookupIdent(ident string) Type {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
