package token

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers + literals
	IDENT  = "IDENT"  // add, foobar, x, y, ...
	INT    = "INT"    // 1343456
	FLOAT  = "FLOAT"  // 12.34
	STRING = "STRING" // "foobar"
	HTML   = "HTML"   // <p>adf</p>

	// Operators
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"
	PERCENT  = "%"

	LT   = "<"
	LTEQ = "<="
	GT   = ">"
	GTEQ = ">="

	EQ     = "=="
	NOT_EQ = "!="
	AND    = "&&"
	OR     = "||"

	// Delimiters

	S_START = "<%"
	C_START = "<%#"
	E_START = "<%="
	E_END   = "%>"

	COMMA     = ","
	SEMICOLON = ";"
	COLON     = ":"

	LPAREN   = "("
	RPAREN   = ")"
	LBRACE   = "{"
	RBRACE   = "}"
	LBRACKET = "["
	RBRACKET = "]"

	// Keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
	FOR      = "FOR"
	IN       = "IN"
)
