package plush

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Render_Function_Call(t *testing.T) {
	r := require.New(t)

	input := `<p><%= f() %></p>`
	s, err := Render(input, NewContextWith(map[string]interface{}{
		"f": func() string {
			return "hi!"
		},
	}))
	r.NoError(err)
	r.Equal("<p>hi!</p>", s)
}

func Test_Render_Unknown_Function_Call(t *testing.T) {
	r := require.New(t)

	input := `<p><%= f() %></p>`
	_, err := Render(input, NewContext())
	r.Error(err)
	r.Contains(err.Error(), "f()")
}

func Test_Render_Function_Call_With_Arg(t *testing.T) {
	r := require.New(t)

	input := `<p><%= f("mark") %></p>`
	s, err := Render(input, NewContextWith(map[string]interface{}{
		"f": func(s string) string {
			return fmt.Sprintf("hi %s!", s)
		},
	}))
	r.NoError(err)
	r.Equal("<p>hi mark!</p>", s)
}

func Test_Render_Function_Call_With_Variable_Arg(t *testing.T) {
	r := require.New(t)

	input := `<p><%= f(name) %></p>`
	s, err := Render(input, NewContextWith(map[string]interface{}{
		"f": func(s string) string {
			return fmt.Sprintf("hi %s!", s)
		},
		"name": "mark",
	}))
	r.NoError(err)
	r.Equal("<p>hi mark!</p>", s)
}

func Test_Render_Function_Call_With_Hash(t *testing.T) {
	r := require.New(t)

	input := `<p><%= f({name: name}) %></p>`
	s, err := Render(input, NewContextWith(map[string]interface{}{
		"f": func(m map[string]interface{}) string {
			return fmt.Sprintf("hi %s!", m["name"])
		},
		"name": "mark",
	}))
	r.NoError(err)
	r.Equal("<p>hi mark!</p>", s)
}

func Test_Render_Function_Call_With_Error(t *testing.T) {
	r := require.New(t)

	input := `<p><%= f() %></p>`
	_, err := Render(input, NewContextWith(map[string]interface{}{
		"f": func() (string, error) {
			return "hi!", errors.New("oops")
		},
	}))
	r.Error(err)
}

func Test_Render_Function_Call_With_Block(t *testing.T) {
	r := require.New(t)

	input := `<p><%= f() { %>hello<% } %></p>`
	s, err := Render(input, NewContextWith(map[string]interface{}{
		"f": func(h HelperContext) string {
			s, _ := h.Block()
			return s
		},
	}))
	r.NoError(err)
	r.Equal("<p>hello</p>", s)
}

type greeter struct{}

func (g greeter) Greet(s string) string {
	return fmt.Sprintf("hi %s!", s)
}

func Test_Render_Function_Call_On_Callee(t *testing.T) {
	r := require.New(t)

	input := `<p><%= g.Greet("mark") %></p>`
	s, err := Render(input, NewContextWith(map[string]interface{}{
		"g": greeter{},
	}))
	r.NoError(err)
	r.Equal(`<p>hi mark!</p>`, s)
}

func Test_Render_Function_Optional_Map(t *testing.T) {
	r := require.New(t)
	input := `<%= foo() %>|<%= bar({a: "A"}) %>`
	s, err := Render(input, NewContextWith(map[string]interface{}{
		"foo": func(opts map[string]interface{}, help HelperContext) string {
			return "foo"
		},
		"bar": func(opts map[string]interface{}) string {
			return opts["a"].(string)
		},
	}))
	r.NoError(err)
	r.Equal("foo|A", s)
}
