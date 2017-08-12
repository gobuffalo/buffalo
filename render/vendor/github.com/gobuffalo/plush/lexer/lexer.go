package lexer

import (
	"regexp"
	"strings"

	"github.com/gobuffalo/plush/token"
)

// Lexer moves through the source input and tokenizes its content
type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
	inside       bool
}

// New Lexer from the input string
func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

// NextToken from the source input
func (l *Lexer) NextToken() token.Token {
	if l.inside {
		return l.nextInsideToken()
	}
	var tok token.Token

	// l.skipWhitespace()
	if l.ch == 0 {
		tok.Literal = ""
		tok.Type = token.EOF
		return tok
	}

	if l.ch == '<' && l.peekChar() == '%' {
		l.inside = true
		return l.nextInsideToken()
	}

	tok.Type = token.HTML
	tok.Literal = l.readHTML()
	return tok
}

func (l *Lexer) nextInsideToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.EQ, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.ASSIGN, l.ch)
		}
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '&':
		if l.peekChar() == '&' {
			l.readChar()
			tok = token.Token{Type: token.AND, Literal: "&&"}
			break
		}
		tok = newToken(token.ILLEGAL, l.ch)
	case '|':
		if l.peekChar() == '|' {
			l.readChar()
			tok = token.Token{Type: token.OR, Literal: "||"}
			break
		}
		tok = newToken(token.ILLEGAL, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.NOT_EQ, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.BANG, l.ch)
		}
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '%':
		if l.peekChar() == '>' {
			l.inside = false
			l.readChar()
			tok = token.Token{Type: token.E_END, Literal: "%>"}
			break
		}
		tok = newToken(token.ILLEGAL, l.ch)
	case '<':
		if l.peekChar() == '%' {
			l.inside = true
			l.readChar()
			switch l.peekChar() {
			case '#':
				l.readChar()
				tok = token.Token{Type: token.C_START, Literal: "<%#"}
			case '=':
				l.readChar()
				tok = token.Token{Type: token.E_START, Literal: "<%="}
			default:
				tok = token.Token{Type: token.S_START, Literal: "<%"}
			}
			break
		}
		if l.peekChar() == '=' {
			l.readChar()
			tok = token.Token{Type: token.LTEQ, Literal: "<="}
			break
		}
		tok = newToken(token.LT, l.ch)
	case '>':
		if l.peekChar() == '=' {
			l.readChar()
			tok = token.Token{Type: token.GTEQ, Literal: ">="}
			break
		}
		tok = newToken(token.GT, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case ':':
		tok = newToken(token.COLON, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString()
	case '[':
		tok = newToken(token.LBRACKET, l.ch)
	case ']':
		tok = newToken(token.RBRACKET, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			if floatX.MatchString(tok.Literal) {
				tok.Type = token.FLOAT
			}
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tok
}

var floatX = regexp.MustCompile(`\d*\.\d*`)

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) || isDigit(l.ch) || l.ch == '.' {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readString() string {
	position := l.position + 1
	for {
		l.readChar()
		// check for quote escapes
		if l.ch == '\\' && l.peekChar() == '"' {
			l.readChar()
			l.readChar()
		}
		if l.ch == '"' {
			break
		}
	}
	s := l.input[position:l.position]
	return strings.Replace(s, "\\\"", "\"", -1)
}

func (l *Lexer) readHTML() string {
	position := l.position

	for l.ch != 0 {
		// allow for expression escaping using \<% foo %>
		if l.ch == '\\' && l.peekChar() == '<' {
			l.readChar()
			l.readChar()
		}
		if l.ch == '<' && l.peekChar() == '%' {
			l.inside = true
			break
		}
		l.readChar()
	}
	return strings.Replace(l.input[position:l.position], "\\<%", "<%", -1)
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_' || ch == '-'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9' || ch == '.'
}

func newToken(tokenType token.Type, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}
