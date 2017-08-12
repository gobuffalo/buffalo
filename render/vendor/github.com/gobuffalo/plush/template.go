package plush

import (
	"github.com/gobuffalo/plush/ast"

	"github.com/gobuffalo/plush/parser"

	"github.com/pkg/errors"
)

// Template represents an input and helpers to be used
// to evaluate and render the input.
type Template struct {
	Input   string
	program *ast.Program
}

// NewTemplate from the input string. Adds all of the
// global helper functions from "Helpers", this function does not
// cache the template.
func NewTemplate(input string) (*Template, error) {
	t := &Template{
		Input: input,
	}
	err := t.Parse()
	if err != nil {
		return t, errors.WithStack(err)
	}
	return t, nil
}

// Parse the template this can be called many times
// as a successful result is cached and is used on subsequent
// uses.
func (t *Template) Parse() error {
	if t.program != nil {
		return nil
	}
	program, err := parser.Parse(t.Input)
	if err != nil {
		return errors.WithStack(err)
	}
	t.program = program
	return nil
}

// Exec the template using the content and return the results
func (t *Template) Exec(ctx *Context) (string, error) {
	err := t.Parse()
	if err != nil {
		return "", err
	}

	// ctx = ctx.New()
	moot.Lock()
	for k, v := range Helpers.helpers {
		ctx.Set(k, v)
	}
	moot.Unlock()
	ev := compiler{
		ctx:     ctx,
		program: t.program,
	}

	s, err := ev.compile()
	return s, err
}

// Clone a template. This is useful for defining helpers on per "instance" of the template.
func (t *Template) Clone() *Template {
	t2 := &Template{
		Input:   t.Input,
		program: t.program,
	}
	return t2
}
