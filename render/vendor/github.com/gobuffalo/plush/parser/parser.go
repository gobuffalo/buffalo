package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gobuffalo/plush/ast"
	"github.com/gobuffalo/plush/lexer"
	"github.com/gobuffalo/plush/token"
)

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

// Parse the string and return an AST or an error
func Parse(s string) (*ast.Program, error) {
	p := newParser(lexer.New(s))
	prog := p.parseProgram()
	if len(p.errors) > 0 {
		return prog, p.errors
	}
	return prog, nil
}

func newParser(l *lexer.Lexer) *parser {
	p := &parser{
		Lexer:  l,
		errors: []string{},
	}

	p.prefixParseFns = make(map[token.Type]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.FLOAT, p.parseFloatLiteral)
	p.registerPrefix(token.STRING, p.parseStringLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.FOR, p.parseForExpression)
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)
	p.registerPrefix(token.LBRACKET, p.parseArrayLiteral)
	p.registerPrefix(token.LBRACE, p.parseHashLiteral)
	p.registerPrefix(token.HTML, p.parseHTMLLiteral)
	p.registerPrefix(token.C_START, p.parseCommentLiteral)
	p.registerPrefix(token.E_END, func() ast.Expression { return nil })

	p.infixParseFns = make(map[token.Type]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.LTEQ, p.parseInfixExpression)
	p.registerInfix(token.GTEQ, p.parseInfixExpression)
	p.registerInfix(token.AND, p.parseInfixExpression)
	p.registerInfix(token.OR, p.parseInfixExpression)

	p.registerInfix(token.LPAREN, p.parseCallExpression)
	p.registerInfix(token.LBRACKET, p.parseIndexExpression)

	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p
}

type parser struct {
	*lexer.Lexer
	errors errSlice

	curToken  token.Token
	peekToken token.Token

	prefixParseFns map[token.Type]prefixParseFn
	infixParseFns  map[token.Type]infixParseFn
}

func (p *parser) parseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if t, ok := stmt.(*ast.ExpressionStatement); ok {
			if _, ok := t.Expression.(*ast.HTMLLiteral); ok {
				program.Statements = append(program.Statements, stmt)
				p.nextToken()
				continue
			}
		}
		if stmt != nil && strings.TrimSpace(stmt.String()) != "" {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

func (p *parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.NextToken()
}

func (p *parser) curTokenIs(t token.Type) bool {
	return p.curToken.Type == t
}

func (p *parser) peekTokenIs(t token.Type) bool {
	return p.peekToken.Type == t
}

func (p *parser) expectPeek(t token.Type) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

func (p *parser) peekError(t token.Type) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *parser) noPrefixParseFnError(t token.Type) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

func (p *parser) parseStatement() ast.Statement {
	// fmt.Println("parseStatement")
	switch p.curToken.Type {
	case token.LET:
		l := p.parseLetStatement()
		return l
	case token.S_START:
		p.nextToken()
		return p.parseStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	case token.E_START:
		return p.parseReturnStatement()
	case token.RBRACE:
		return nil
	case token.EOF:
		return nil
	default:
		return p.parseExpressionStatement()
	}
}

func (p *parser) parseReturnStatement() *ast.ReturnStatement {
	// fmt.Println("parseReturnStatement")
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	stmt.ReturnValue = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *parser) parseLetStatement() *ast.LetStatement {
	// fmt.Println("parseLetStatement")
	stmt := &ast.LetStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *parser) parseExpressionStatement() *ast.ExpressionStatement {
	// fmt.Println("parseExpressionStatement")
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if p.curTokenIs(token.LET) {
		return nil
	}
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()

		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *parser) parseIdentifier() ast.Expression {
	// fmt.Println("parseIdentifier")
	id := &ast.Identifier{Token: p.curToken}
	ss := strings.Split(p.curToken.Literal, ".")
	id.Value = ss[0]

	for i := 1; i < len(ss); i++ {
		s := ss[i]
		id = &ast.Identifier{Token: p.curToken, Value: s, Callee: id}
	}
	return id
}

func (p *parser) parseIntegerLiteral() ast.Expression {
	// fmt.Println("parseIntegerLiteral")
	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.Atoi(p.curToken.Literal)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value

	return lit
}

func (p *parser) parseFloatLiteral() ast.Expression {
	// fmt.Println("parseFloatLiteral")
	lit := &ast.FloatLiteral{Token: p.curToken}

	value, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as float", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value

	return lit
}

func (p *parser) parseStringLiteral() ast.Expression {
	// fmt.Println("parseStringLiteral")
	return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *parser) parseCommentLiteral() ast.Expression {
	// fmt.Println("parseCommentLiteral")
	for p.curToken.Type != token.E_END {
		p.nextToken()
	}
	return &ast.StringLiteral{Token: p.curToken, Value: ""}
}

func (p *parser) parseHTMLLiteral() ast.Expression {
	// fmt.Println("parseHTMLLiteral")
	return &ast.HTMLLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *parser) parsePrefixExpression() ast.Expression {
	// fmt.Println("parsePrefixExpression")
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()

	expression.Right = p.parseExpression(PREFIX)

	return expression
}

func (p *parser) parseInfixExpression(left ast.Expression) ast.Expression {
	// fmt.Println("parseInfixExpression")
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *parser) parseBoolean() ast.Expression {
	// fmt.Println("parseBoolean")
	return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
}

func (p *parser) parseGroupedExpression() ast.Expression {
	// fmt.Println("parseGroupedExpression")
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return exp
}

func (p *parser) parseForExpression() ast.Expression {
	// fmt.Println("parseForExpression")
	expression := &ast.ForExpression{
		Token:     p.curToken,
		KeyName:   "_",
		ValueName: "@value",
	}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	s := []string{}
	for !p.curTokenIs(token.RPAREN) {
		if p.curTokenIs(token.IDENT) {
			s = append(s, p.curToken.Literal)
		}
		p.nextToken()
	}

	switch len(s) {
	case 1:
		expression.ValueName = s[0]
	case 2:
		expression.KeyName = s[0]
		expression.ValueName = s[1]
	}

	p.nextToken()

	if !p.curTokenIs(token.IN) {
		return nil
	}
	p.nextToken()
	expression.Iterable = p.parseExpression(LOWEST)
	if ce, ok := expression.Iterable.(*ast.CallExpression); ok {
		if ce.Block != nil {
			expression.Block = ce.Block
			ce.Block = nil
			return expression
		}
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	expression.Block = p.parseBlockStatement()

	if p.curTokenIs(token.RBRACE) {
		p.nextToken()
	}

	return expression
}

func (p *parser) parseIfExpression() ast.Expression {
	// fmt.Println("parseIfExpression")
	expression := &ast.IfExpression{Token: p.curToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	expression.Block = p.parseBlockStatement()

	if p.peekTokenIs(token.ELSE) {
		p.nextToken()

		if !p.expectPeek(token.LBRACE) {
			return nil
		}

		expression.ElseBlock = p.parseBlockStatement()
	}

	return expression
}

func (p *parser) parseBlockStatement() *ast.BlockStatement {
	// fmt.Println("parseBlockStatement")
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}

	p.nextToken()

	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		if p.curTokenIs(token.S_START) || p.curTokenIs(token.E_END) {
			p.nextToken()
			continue
		}
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}

func (p *parser) parseFunctionLiteral() ast.Expression {
	// fmt.Println("parseFunctionLiteral")
	lit := &ast.FunctionLiteral{Token: p.curToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	lit.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	lit.Block = p.parseBlockStatement()

	return lit
}

func (p *parser) parseFunctionParameters() []*ast.Identifier {
	// fmt.Println("parseFunctionParameters")
	identifiers := []*ast.Identifier{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return identifiers
	}

	p.nextToken()

	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	identifiers = append(identifiers, ident)

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		identifiers = append(identifiers, ident)
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return identifiers
}

func (p *parser) parseCallExpression(function ast.Expression) ast.Expression {
	// fmt.Println("parseCallExpression")
	exp := &ast.CallExpression{
		Token:    p.curToken,
		Function: function,
	}
	ss := strings.Split(function.String(), ".")
	if len(ss) > 1 {
		exp.Callee = &ast.Identifier{
			Token: token.Token{Type: token.IDENT, Literal: ss[0]},
			Value: ss[0],
		}
		for i := 1; i < len(ss)-1; i++ {
			c := &ast.Identifier{
				Token:  token.Token{Type: token.IDENT, Literal: ss[i]},
				Value:  ss[i],
				Callee: exp.Callee.(*ast.Identifier),
			}
			exp.Callee = c
		}
		exp.Function = &ast.Identifier{
			Token:  token.Token{Type: token.IDENT, Literal: ss[len(ss)-1]},
			Value:  ss[len(ss)-1],
			Callee: exp.Callee.(*ast.Identifier),
		}
	}
	exp.Arguments = p.parseExpressionList(token.RPAREN)

	if p.peekTokenIs(token.LBRACE) {
		p.nextToken()

		exp.Block = p.parseBlockStatement()
	}
	return exp
}

func (p *parser) parseExpressionList(end token.Type) []ast.Expression {
	// fmt.Println("parseExpressionList")
	list := []ast.Expression{}

	if p.peekTokenIs(end) {
		p.nextToken()
		return list
	}

	p.nextToken()
	list = append(list, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		list = append(list, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(end) {
		return nil
	}

	return list
}

func (p *parser) parseArrayLiteral() ast.Expression {
	// fmt.Println("parseArrayLiteral")
	array := &ast.ArrayLiteral{Token: p.curToken}

	array.Elements = p.parseExpressionList(token.RBRACKET)

	return array
}

func (p *parser) parseIndexExpression(left ast.Expression) ast.Expression {
	// fmt.Println("parseIndexExpression")
	exp := &ast.IndexExpression{Token: p.curToken, Left: left}

	p.nextToken()
	exp.Index = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RBRACKET) {
		return nil
	}

	return exp
}

func (p *parser) parseHashLiteral() ast.Expression {
	// fmt.Println("parseHashLiteral")
	hash := &ast.HashLiteral{Token: p.curToken}
	hash.Pairs = make(map[ast.Expression]ast.Expression)

	for !p.peekTokenIs(token.RBRACE) {
		p.nextToken()
		key := p.parseExpression(LOWEST)

		if !p.expectPeek(token.COLON) {
			return nil
		}

		p.nextToken()
		value := p.parseExpression(LOWEST)

		hash.Pairs[key] = value

		if !p.peekTokenIs(token.RBRACE) && !p.expectPeek(token.COMMA) {
			return nil
		}
	}

	if !p.expectPeek(token.RBRACE) {
		return nil
	}

	return hash
}

func (p *parser) registerPrefix(tokenType token.Type, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *parser) registerInfix(tokenType token.Type, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}
