package lexer

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/gobuffalo/plush/token"
)

func Test_NextToken_Simple(t *testing.T) {
	r := require.New(t)
	input := `<%= 1 %>`
	tests := []struct {
		tokenType    token.Type
		tokenLiteral string
	}{
		{token.E_START, "<%="},
		{token.INT, "1"},
		{token.E_END, "%>"},
	}

	l := New(input)
	for _, tt := range tests {
		tok := l.NextToken()
		r.Equal(tt.tokenType, tok.Type)
		r.Equal(tt.tokenLiteral, tok.Literal)
	}
}

func Test_EscapeStringQuote(t *testing.T) {
	r := require.New(t)
	input := `<%= "mark \"cool\" bates" %>`
	tests := []struct {
		tokenType    token.Type
		tokenLiteral string
	}{
		{token.E_START, "<%="},
		{token.STRING, `mark "cool" bates`},
		{token.E_END, "%>"},
	}

	l := New(input)
	for _, tt := range tests {
		tok := l.NextToken()
		r.Equal(tt.tokenType, tok.Type)
		r.Equal(tt.tokenLiteral, tok.Literal)
	}
}

func Test_EscapeExpression(t *testing.T) {
	r := require.New(t)
	input := `<p>\<%= 1 %></p>`
	tests := []struct {
		tokenType    token.Type
		tokenLiteral string
	}{
		{token.HTML, `<p><%= 1 %></p>`},
	}

	l := New(input)
	for _, tt := range tests {
		tok := l.NextToken()
		r.Equal(tt.tokenType, tok.Type)
		r.Equal(tt.tokenLiteral, tok.Literal)
	}
}

func Test_NextToken_WithHTML(t *testing.T) {
	r := require.New(t)
	input := `<p class="foo"><%= 1 %></p>`
	tests := []struct {
		tokenType    token.Type
		tokenLiteral string
	}{
		{token.HTML, `<p class="foo">`},
		{token.E_START, "<%="},
		{token.INT, "1"},
		{token.E_END, "%>"},
		{token.HTML, `</p>`},
	}

	l := New(input)
	for _, tt := range tests {
		tok := l.NextToken()
		r.Equal(tt.tokenType, tok.Type)
		r.Equal(tt.tokenLiteral, tok.Literal)
	}
}
func Test_NextToken_Complete(t *testing.T) {
	r := require.New(t)
	input := `<% let five = 5;
let ten = 10;

let add = fn(x, y) {
  x + y;
};

let result = add(five, ten);
!-/*5;
5 < 10 > 5;

if (5 < 10) {
	return true;
} else {
	return false;
}

10 == 10;
10 != 9;
"foobar"
"foo bar"
[1, 2];
{"foo": "bar"}
let fl = 1.23 %>
<%= 1 %>
<%# 2 %>
<% 3 %>
<% for (i, v) in myArray {
}
a && b
c || d
for (x) in range(1,3){return x}
myvar1
my-helper()
%>
`

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.S_START, "<%"},
		{token.LET, "let"},
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "fn"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "result"},
		{token.ASSIGN, "="},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "five"},
		{token.COMMA, ","},
		{token.IDENT, "ten"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.BANG, "!"},
		{token.MINUS, "-"},
		{token.SLASH, "/"},
		{token.ASTERISK, "*"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.GT, ">"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.TRUE, "true"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.ELSE, "else"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.FALSE, "false"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.INT, "10"},
		{token.EQ, "=="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.INT, "10"},
		{token.NOT_EQ, "!="},
		{token.INT, "9"},
		{token.SEMICOLON, ";"},
		{token.STRING, "foobar"},
		{token.STRING, "foo bar"},
		{token.LBRACKET, "["},
		{token.INT, "1"},
		{token.COMMA, ","},
		{token.INT, "2"},
		{token.RBRACKET, "]"},
		{token.SEMICOLON, ";"},
		{token.LBRACE, "{"},
		{token.STRING, "foo"},
		{token.COLON, ":"},
		{token.STRING, "bar"},
		{token.RBRACE, "}"},
		{token.LET, "let"},
		{token.IDENT, "fl"},
		{token.ASSIGN, "="},
		{token.FLOAT, "1.23"},
		{token.E_END, "%>"},
		{token.HTML, "\n"},
		{token.E_START, "<%="},
		{token.INT, "1"},
		{token.E_END, "%>"},
		{token.HTML, "\n"},
		{token.C_START, "<%#"},
		{token.INT, "2"},
		{token.E_END, "%>"},
		{token.HTML, "\n"},
		{token.S_START, "<%"},
		{token.INT, "3"},
		{token.E_END, "%>"},
		{token.HTML, "\n"},
		{token.S_START, "<%"},
		{token.FOR, "for"},
		{token.LPAREN, "("},
		{token.IDENT, "i"},
		{token.COMMA, ","},
		{token.IDENT, "v"},
		{token.RPAREN, ")"},
		{token.IN, "in"},
		{token.IDENT, "myArray"},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.IDENT, "a"},
		{token.AND, "&&"},
		{token.IDENT, "b"},
		{token.IDENT, "c"},
		{token.OR, "||"},
		{token.IDENT, "d"},
		{token.FOR, "for"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.RPAREN, ")"},
		{token.IN, "in"},
		{token.IDENT, "range"},
		{token.LPAREN, "("},
		{token.INT, "1"},
		{token.COMMA, ","},
		{token.INT, "3"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.IDENT, "x"},
		{token.RBRACE, "}"},
		{token.IDENT, "myvar1"},
		{token.IDENT, "my-helper"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.E_END, "%>"},
		{token.HTML, "\n"},
		{token.EOF, ""},
	}

	l := New(input)

	for _, tt := range tests {
		tok := l.NextToken()

		r.Equal(tt.expectedLiteral, tok.Literal)
		r.Equal(tt.expectedType, tok.Type)
	}
}
