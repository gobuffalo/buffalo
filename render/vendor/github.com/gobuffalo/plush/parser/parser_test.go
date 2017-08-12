package parser

import (
	"fmt"
	"testing"

	"github.com/gobuffalo/plush/ast"

	"github.com/stretchr/testify/require"
)

func Test_LetStatements(t *testing.T) {
	r := require.New(t)
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"<% let x = 5; %>", "x", 5},
		{"<% let y = true; %>", "y", true},
		{"<% let foobar = y;%>", "foobar", "y"},
	}

	for _, tt := range tests {
		program, err := Parse(tt.input)
		r.NoError(err)

		r.Len(program.Statements, 1)
		stmt := program.Statements[0]

		r.Equal("let", stmt.TokenLiteral())

		letStmt := stmt.(*ast.LetStatement)

		r.Equal(tt.expectedIdentifier, letStmt.Name.Value)
		r.Equal(tt.expectedIdentifier, letStmt.Name.TokenLiteral())

		val := stmt.(*ast.LetStatement).Value
		r.True(testLiteralExpression(t, val, tt.expectedValue))
	}
}

func Test_ReturnStatements(t *testing.T) {
	r := require.New(t)
	tests := []struct {
		input         string
		expectedValue interface{}
	}{
		{"return 5;", 5},
		{"return true;", true},
		{"return foobar;", "foobar"},
	}

	for _, tt := range tests {
		program, err := Parse("<%" + tt.input + "%>")
		r.NoError(err)

		r.Len(program.Statements, 1)

		stmt := program.Statements[0]
		returnStmt := stmt.(*ast.ReturnStatement)
		r.Equal("return", returnStmt.TokenLiteral())
		r.True(testLiteralExpression(t, returnStmt.ReturnValue, tt.expectedValue))
	}
}

func Test_IdentifierExpression(t *testing.T) {
	r := require.New(t)
	input := "<% foobar; %>"

	program, err := Parse(input)
	r.NoError(err)

	r.Len(program.Statements, 1)
	stmt := program.Statements[0].(*ast.ExpressionStatement)
	ident := stmt.Expression.(*ast.Identifier)
	r.Equal("foobar", ident.Value)
	r.Equal("foobar", ident.TokenLiteral())
}

func Test_IntegerLiteralExpression(t *testing.T) {
	r := require.New(t)
	input := "<% 5; %>"

	program, err := Parse(input)
	r.NoError(err)

	r.Len(program.Statements, 1)
	stmt := program.Statements[0].(*ast.ExpressionStatement)
	literal := stmt.Expression.(*ast.IntegerLiteral)
	r.Equal(5, literal.Value)
	r.Equal("5", literal.TokenLiteral())
}

func Test_FloatLiteralExpression(t *testing.T) {
	r := require.New(t)
	input := "<% 1.23 %>"

	program, err := Parse(input)
	r.NoError(err)

	r.Len(program.Statements, 1)
	stmt := program.Statements[0].(*ast.ExpressionStatement)

	literal := stmt.Expression.(*ast.FloatLiteral)
	r.Equal(1.23, literal.Value)
	r.Equal("1.23", literal.TokenLiteral())
}

func Test_PrefixExpressions(t *testing.T) {
	r := require.New(t)
	prefixTests := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
		{"!foobar;", "!", "foobar"},
		{"-foobar;", "-", "foobar"},
		{"!true;", "!", true},
		{"!false;", "!", false},
	}

	for _, tt := range prefixTests {
		program, err := Parse("<%" + tt.input + "%>")
		r.NoError(err)

		r.Len(program.Statements, 1)

		stmt := program.Statements[0].(*ast.ExpressionStatement)

		exp := stmt.Expression.(*ast.PrefixExpression)
		r.Equal(tt.operator, exp.Operator)

		r.True(testLiteralExpression(t, exp.Right, tt.value))
	}
}

func Test_InfixExpressions(t *testing.T) {
	r := require.New(t)
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"foobar + barfoo;", "foobar", "+", "barfoo"},
		{"foobar - barfoo;", "foobar", "-", "barfoo"},
		{"foobar * barfoo;", "foobar", "*", "barfoo"},
		{"foobar / barfoo;", "foobar", "/", "barfoo"},
		{"foobar > barfoo;", "foobar", ">", "barfoo"},
		{"foobar < barfoo;", "foobar", "<", "barfoo"},
		{"foobar == barfoo;", "foobar", "==", "barfoo"},
		{"foobar != barfoo;", "foobar", "!=", "barfoo"},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}

	for _, tt := range infixTests {
		program, err := Parse("<%" + tt.input + "%>")
		r.NoError(err)

		r.Len(program.Statements, 1)
		stmt := program.Statements[0].(*ast.ExpressionStatement)
		r.True(testInfixExpression(t, stmt.Expression, tt.leftValue, tt.operator, tt.rightValue))
	}
}

func Test_OperatorPrecedence(t *testing.T) {
	r := require.New(t)
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"(5 + 5) * 2 * (5 + 5)",
			"(((5 + 5) * 2) * (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
		{
			"a + add(b * c) + d",
			"((a + add((b * c))) + d)",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		},
		{
			"add(a + b + c * d / f + g)",
			"add((((a + b) + ((c * d) / f)) + g))",
		},
		{
			"a * [1, 2, 3, 4][b * c] * d",
			"((a * ([1, 2, 3, 4][(b * c)])) * d)",
		},
		{
			"add(a * b[2], b[1], 2 * [1, 2][1])",
			"add((a * (b[2])), (b[1]), (2 * ([1, 2][1])))",
		},
	}

	for _, tt := range tests {
		program, err := Parse("<%" + tt.input + "%>")
		r.NoError(err)

		r.Equal(tt.expected, program.String())
	}
}

func Test_BooleanExpression(t *testing.T) {
	r := require.New(t)
	tests := []struct {
		input           string
		expectedBoolean bool
	}{
		{"true;", true},
		{"false;", false},
	}

	for _, tt := range tests {
		program, err := Parse("<%" + tt.input + "%>")
		r.NoError(err)

		r.Len(program.Statements, 1)

		stmt := program.Statements[0].(*ast.ExpressionStatement)

		boolean := stmt.Expression.(*ast.Boolean)
		r.Equal(tt.expectedBoolean, boolean.Value)
	}
}

func Test_IfExpression(t *testing.T) {
	r := require.New(t)
	input := `<% if (x < y) { x } %>`

	program, err := Parse(input)
	r.NoError(err)

	r.Len(program.Statements, 1)

	stmt := program.Statements[0].(*ast.ExpressionStatement)

	exp := stmt.Expression.(*ast.IfExpression)

	r.True(testInfixExpression(t, exp.Condition, "x", "<", "y"))

	r.Len(exp.Block.Statements, 1)

	consequence := exp.Block.Statements[0].(*ast.ExpressionStatement)

	r.True(testIdentifier(t, consequence.Expression, "x"))
	r.Nil(exp.ElseBlock)
}

func Test_IfExpression_HTML(t *testing.T) {
	r := require.New(t)
	input := `<p><% if (x < y) { %><%= x %><% } %></p>`

	program, err := Parse(input)
	r.NoError(err)

	r.Len(program.Statements, 3)

	es := program.Statements[0].(*ast.ExpressionStatement)
	h := es.Expression.(*ast.HTMLLiteral)
	r.Equal("<p>", h.Value)

	es = program.Statements[1].(*ast.ExpressionStatement)
	ifs := es.Expression.(*ast.IfExpression)

	r.True(testInfixExpression(t, ifs.Condition, "x", "<", "y"))

	r.Len(ifs.Block.Statements, 1)

	ret := ifs.Block.Statements[0].(*ast.ReturnStatement)

	r.Equal("x", ret.ReturnValue.String())

	r.Nil(ifs.ElseBlock)
}

func Test_IfElseExpression(t *testing.T) {
	r := require.New(t)
	input := `<% if (x < y) { x } else { y } %>`

	program, err := Parse(input)
	r.NoError(err)

	r.Len(program.Statements, 1)

	stmt := program.Statements[0].(*ast.ExpressionStatement)

	exp := stmt.Expression.(*ast.IfExpression)

	r.True(testInfixExpression(t, exp.Condition, "x", "<", "y"))

	r.Len(exp.Block.Statements, 1)

	consequence := exp.Block.Statements[0].(*ast.ExpressionStatement)

	r.True(testIdentifier(t, consequence.Expression, "x"))

	r.Len(exp.ElseBlock.Statements, 1)

	alternative := exp.ElseBlock.Statements[0].(*ast.ExpressionStatement)

	r.True(testIdentifier(t, alternative.Expression, "y"))
}

func Test_FunctionLiteralParsing(t *testing.T) {
	r := require.New(t)
	input := `<% fn(x, y) { x + y; } %>`

	program, err := Parse(input)
	r.NoError(err)

	r.Len(program.Statements, 1)

	stmt := program.Statements[0].(*ast.ExpressionStatement)

	function := stmt.Expression.(*ast.FunctionLiteral)

	r.Len(function.Parameters, 2)

	r.True(testLiteralExpression(t, function.Parameters[0], "x"))
	r.True(testLiteralExpression(t, function.Parameters[1], "y"))

	r.Len(function.Block.Statements, 1)

	bodyStmt := function.Block.Statements[0].(*ast.ExpressionStatement)

	r.True(testInfixExpression(t, bodyStmt.Expression, "x", "+", "y"))
}

func Test_FunctionParameterParsing(t *testing.T) {
	r := require.New(t)
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{input: "fn() {};", expectedParams: []string{}},
		{input: "fn(x) {};", expectedParams: []string{"x"}},
		{input: "fn(x, y, z) {};", expectedParams: []string{"x", "y", "z"}},
	}

	for _, tt := range tests {
		program, err := Parse("<%" + tt.input + "%>")
		r.NoError(err)

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		function := stmt.Expression.(*ast.FunctionLiteral)

		r.Equal(len(function.Parameters), len(tt.expectedParams))
		for i, ident := range tt.expectedParams {
			r.True(testLiteralExpression(t, function.Parameters[i], ident))
		}
	}
}

func Test_CallExpression(t *testing.T) {
	r := require.New(t)
	input := "<% add(1, 2 * 3, 4 + 5); %>"

	program, err := Parse(input)
	r.NoError(err)

	r.Len(program.Statements, 1)

	stmt := program.Statements[0].(*ast.ExpressionStatement)

	exp := stmt.Expression.(*ast.CallExpression)

	r.True(testIdentifier(t, exp.Function, "add"))

	r.Len(exp.Arguments, 3)

	r.True(testLiteralExpression(t, exp.Arguments[0], 1))
	r.True(testInfixExpression(t, exp.Arguments[1], 2, "*", 3))
	r.True(testInfixExpression(t, exp.Arguments[2], 4, "+", 5))
}

func Test_CallExpressionParameter(t *testing.T) {
	r := require.New(t)
	tests := []struct {
		input         string
		expectedIdent string
		expectedArgs  []string
	}{
		{
			input:         "add();",
			expectedIdent: "add",
			expectedArgs:  []string{},
		},
		{
			input:         "add(1);",
			expectedIdent: "add",
			expectedArgs:  []string{"1"},
		},
		{
			input:         "add(1, 2 * 3, 4 + 5);",
			expectedIdent: "add",
			expectedArgs:  []string{"1", "(2 * 3)", "(4 + 5)"},
		},
	}

	for _, tt := range tests {
		program, err := Parse("<%" + tt.input + "%>")
		r.NoError(err)

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		exp := stmt.Expression.(*ast.CallExpression)

		r.True(testIdentifier(t, exp.Function, tt.expectedIdent))

		r.Equal(len(exp.Arguments), len(tt.expectedArgs))

		for i, arg := range tt.expectedArgs {
			r.Equal(arg, exp.Arguments[i].String())
		}
	}
}

func Test_CallExpressionParsing_WithCallee(t *testing.T) {
	r := require.New(t)
	input := `<%= g.Greet("mark"); %>`

	program, err := Parse(input)
	r.NoError(err)

	r.Len(program.Statements, 1)

	stmt := program.Statements[0].(*ast.ReturnStatement)

	exp := stmt.ReturnValue.(*ast.CallExpression)

	ident := exp.Function.(*ast.Identifier)
	r.Equal("Greet", ident.Value)

	r.Len(exp.Arguments, 1)
	r.Equal(exp.Arguments[0].String(), "mark")
}

func Test_CallExpressionParsing_WithMultipleCallees(t *testing.T) {
	r := require.New(t)
	input := `<%= g.Foo.Greet("mark"); %>`

	program, err := Parse(input)
	r.NoError(err)

	r.Len(program.Statements, 1)

	stmt := program.Statements[0].(*ast.ReturnStatement)

	exp := stmt.ReturnValue.(*ast.CallExpression)

	ident := exp.Function.(*ast.Identifier)
	r.Equal("Greet", ident.Value)

	r.Len(exp.Arguments, 1)
	r.Equal(exp.Arguments[0].String(), "mark")
}

func Test_CallExpressionParsing_WithBlock(t *testing.T) {
	r := require.New(t)
	input := `<p><%= foo() { %>hi<% } %></p>`

	program, err := Parse(input)
	r.NoError(err)

	r.Len(program.Statements, 3)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	html := stmt.Expression.(*ast.HTMLLiteral)
	r.Equal("<p>", html.Value)

	rstmt := program.Statements[1].(*ast.ReturnStatement)
	exp := rstmt.ReturnValue.(*ast.CallExpression)

	ident := exp.Function.(*ast.Identifier)
	r.Equal("foo", ident.Value)

	r.Len(exp.Arguments, 0)
	r.NotNil(exp.Block)
	r.Equal("hi", exp.Block.String())
	r.Nil(exp.Callee)
}

func Test_StringLiteralExpression(t *testing.T) {
	r := require.New(t)
	input := `<% "hello world"; %>`

	program, err := Parse(input)
	r.NoError(err)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	literal := stmt.Expression.(*ast.StringLiteral)

	r.Equal("hello world", literal.Value)
}

func Test_EmptyArrayLiterals(t *testing.T) {
	r := require.New(t)
	input := "<% [] %>"

	program, err := Parse(input)
	r.NoError(err)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	array := stmt.Expression.(*ast.ArrayLiteral)

	r.Len(array.Elements, 0)
}

func Test_ArrayLiterals(t *testing.T) {
	r := require.New(t)
	input := "<% [1, 2 * 2, 3 + 3] %>"

	program, err := Parse(input)
	r.NoError(err)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	array := stmt.Expression.(*ast.ArrayLiteral)

	r.Len(array.Elements, 3)

	r.True(testIntegerLiteral(t, array.Elements[0], 1))
	r.True(testInfixExpression(t, array.Elements[1], 2, "*", 2))
	r.True(testInfixExpression(t, array.Elements[2], 3, "+", 3))
}

func Test_IndexExpressions(t *testing.T) {
	r := require.New(t)
	input := "<% myArray[1 + 1] %>"

	program, err := Parse(input)
	r.NoError(err)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	indexExp := stmt.Expression.(*ast.IndexExpression)

	r.True(testIdentifier(t, indexExp.Left, "myArray"))

	r.True(testInfixExpression(t, indexExp.Index, 1, "+", 1))
}

func Test_EmptyHashLiteral(t *testing.T) {
	r := require.New(t)
	input := "<% {} %>"

	program, err := Parse(input)
	r.NoError(err)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash := stmt.Expression.(*ast.HashLiteral)

	r.Len(hash.Pairs, 0)
}

func Test_HashLiteralsStringKeys(t *testing.T) {
	r := require.New(t)
	input := `<% {"one": 1, "two": 2, "three": 3} %>`

	program, err := Parse(input)
	r.NoError(err)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash := stmt.Expression.(*ast.HashLiteral)

	expected := map[string]int{
		"one":   1,
		"two":   2,
		"three": 3,
	}

	r.Equal(len(hash.Pairs), len(expected))

	for key, value := range hash.Pairs {
		literal := key.(*ast.StringLiteral)

		expectedValue := expected[literal.String()]
		r.True(testIntegerLiteral(t, value, expectedValue))
	}
}

func Test_HashLiteralsBooleanKeys(t *testing.T) {
	r := require.New(t)
	input := `<%{true: 1, false: 2}%>`

	program, err := Parse(input)
	r.NoError(err)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash := stmt.Expression.(*ast.HashLiteral)

	expected := map[string]int{
		"true":  1,
		"false": 2,
	}

	r.Equal(len(hash.Pairs), len(expected))

	for key, value := range hash.Pairs {
		boolean := key.(*ast.Boolean)

		expectedValue := expected[boolean.String()]
		r.True(testIntegerLiteral(t, value, expectedValue))
	}
}

func Test_HashLiteralsIntegerKeys(t *testing.T) {
	r := require.New(t)
	input := `<% {1: 1, 2: 2, 3: 3} %>`

	program, err := Parse(input)
	r.NoError(err)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash := stmt.Expression.(*ast.HashLiteral)

	expected := map[string]int{
		"1": 1,
		"2": 2,
		"3": 3,
	}

	r.Equal(len(hash.Pairs), len(expected))

	for key, value := range hash.Pairs {
		integer := key.(*ast.IntegerLiteral)

		expectedValue := expected[integer.String()]

		r.True(testIntegerLiteral(t, value, expectedValue))
	}
}

func Test_HashLiteralsWithExpressions(t *testing.T) {
	r := require.New(t)
	input := `<% {"one": 0 + 1, "two": 10 - 8, "three": 15 / 5} %>`

	program, err := Parse(input)
	r.NoError(err)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash := stmt.Expression.(*ast.HashLiteral)

	r.Len(hash.Pairs, 3)

	tests := map[string]func(ast.Expression){
		"one": func(e ast.Expression) {
			testInfixExpression(t, e, 0, "+", 1)
		},
		"two": func(e ast.Expression) {
			testInfixExpression(t, e, 10, "-", 8)
		},
		"three": func(e ast.Expression) {
			testInfixExpression(t, e, 15, "/", 5)
		},
	}

	for key, value := range hash.Pairs {
		literal := key.(*ast.StringLiteral)

		testFunc := tests[literal.String()]
		testFunc(value)
	}
}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{}, operator string, right interface{}) bool {
	r := require.New(t)
	opExp := exp.(*ast.InfixExpression)

	r.True(testLiteralExpression(t, opExp.Left, left))

	r.Equal(operator, opExp.Operator)

	r.True(testLiteralExpression(t, opExp.Right, right))
	return true
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, v)
	case int64:
		return testIntegerLiteral(t, exp, int(v))
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)
	}
	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int) bool {
	r := require.New(t)
	integ := il.(*ast.IntegerLiteral)

	r.Equal(value, integ.Value)
	r.Equal(fmt.Sprint(value), integ.TokenLiteral())

	return true
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	r := require.New(t)
	ident := exp.(*ast.Identifier)

	r.Equal(value, ident.Value)
	r.Equal(value, ident.TokenLiteral())

	return true
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
	r := require.New(t)
	bo := exp.(*ast.Boolean)

	r.Equal(value, bo.Value)
	r.Equal(fmt.Sprint(value), bo.TokenLiteral())

	return true
}

func Test_ForExpression(t *testing.T) {
	r := require.New(t)
	input := `<% for (k,v) in myArray { v } %>`

	program, err := Parse(input)
	r.NoError(err)

	r.Len(program.Statements, 1)

	stmt := program.Statements[0].(*ast.ExpressionStatement)

	exp := stmt.Expression.(*ast.ForExpression)

	r.Equal("k", exp.KeyName)
	r.Equal("v", exp.ValueName)
	r.Equal("myArray", exp.Iterable.String())

	r.Len(exp.Block.Statements, 1)

	consequence := exp.Block.Statements[0].(*ast.ExpressionStatement)

	r.True(testIdentifier(t, consequence.Expression, "v"))
}

func Test_ForExpression_Split(t *testing.T) {
	r := require.New(t)
	input := `<% for (k,v) in anArray { %>
	<p><%= v %></p>
	<% } %>`

	program, err := Parse(input)
	r.NoError(err)

	r.Len(program.Statements, 1)

	stmt := program.Statements[0].(*ast.ExpressionStatement)

	exp := stmt.Expression.(*ast.ForExpression)

	r.Equal("k", exp.KeyName)
	r.Equal("v", exp.ValueName)
	r.Equal("anArray", exp.Iterable.String())
	r.Len(exp.Block.Statements, 3)
}

func Test_ForExpression_Func(t *testing.T) {
	r := require.New(t)
	input := `<% for (k,v) in range(1,3) { %>
	<p><%= v %></p>
	<% } %>`

	program, err := Parse(input)
	r.NoError(err)

	r.Len(program.Statements, 1)

	stmt := program.Statements[0].(*ast.ExpressionStatement)

	exp := stmt.Expression.(*ast.ForExpression)

	r.Equal("k", exp.KeyName)
	r.Equal("v", exp.ValueName)
	r.Equal("range(1, 3)", exp.Iterable.String())
	r.Len(exp.Block.Statements, 3)
}

func Test_AndOrInfixExpressions(t *testing.T) {
	r := require.New(t)
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"foobar && barfoo;", "foobar", "&&", "barfoo"},
		{"foobar || barfoo;", "foobar", "||", "barfoo"},
		{"true && true", "true", "&&", "true"},
		{"true || false", "true", "||", "false"},
	}

	for _, tt := range infixTests {
		program, err := Parse("<% " + tt.input + "%>")
		r.NoError(err)

		r.Len(program.Statements, 1)
		stmt := program.Statements[0].(*ast.ExpressionStatement)
		ins := stmt.Expression.(*ast.InfixExpression)
		r.Equal(ins.Right.String(), tt.rightValue)
	}
}
