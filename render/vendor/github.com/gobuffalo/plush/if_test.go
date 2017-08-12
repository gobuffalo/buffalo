package plush

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Render_If(t *testing.T) {
	r := require.New(t)
	input := `<% if (true) { return "hi"} %>`
	s, err := Render(input, NewContext())
	r.NoError(err)
	r.Equal("", s)
}

func Test_Render_If_Return(t *testing.T) {
	r := require.New(t)
	input := `<%= if (true) { return "hi"} %>`
	s, err := Render(input, NewContext())
	r.NoError(err)
	r.Equal("hi", s)
}

func Test_Render_If_Return_HTML(t *testing.T) {
	r := require.New(t)
	input := `<%= if (true) { %>hi<%} %>`
	s, err := Render(input, NewContext())
	r.NoError(err)
	r.Equal("hi", s)
}

func Test_Render_If_And(t *testing.T) {
	r := require.New(t)
	input := `<%= if (false && true) { %> hi <%} %>`
	s, err := Render(input, NewContext())
	r.NoError(err)
	r.Equal("", s)
}

func Test_Render_If_Or(t *testing.T) {
	r := require.New(t)
	input := `<%= if (false || true) { %>hi<%} %>`
	s, err := Render(input, NewContext())
	r.NoError(err)
	r.Equal("hi", s)
}

func Test_Render_If_Nil(t *testing.T) {
	r := require.New(t)
	input := `<%= if (names && len(names) >= 1) { %>hi<%} %>`
	s, err := Render(input, NewContext())
	r.NoError(err)
	r.Equal("", s)
}

func Test_Render_If_Else_Return(t *testing.T) {
	r := require.New(t)
	input := `<p><%= if (false) { return "hi"} else { return "bye"} %></p>`
	s, err := Render(input, NewContext())
	r.NoError(err)
	r.Equal("<p>bye</p>", s)
}

func Test_Render_If_LessThan(t *testing.T) {
	r := require.New(t)
	input := `<p><%= if (1 < 2) { return "hi"} else { return "bye"} %></p>`
	s, err := Render(input, NewContext())
	r.NoError(err)
	r.Equal("<p>hi</p>", s)
}

func Test_Render_If_BangFalse(t *testing.T) {
	r := require.New(t)
	input := `<p><%= if (!false) { return "hi"} else { return "bye"} %></p>`
	s, err := Render(input, NewContext())
	r.NoError(err)
	r.Equal("<p>hi</p>", s)
}

func Test_Render_If_NotEq(t *testing.T) {
	r := require.New(t)
	input := `<p><%= if (1 != 2) { return "hi"} else { return "bye"} %></p>`
	s, err := Render(input, NewContext())
	r.NoError(err)
	r.Equal("<p>hi</p>", s)
}

func Test_Render_If_GtEq(t *testing.T) {
	r := require.New(t)
	input := `<p><%= if (1 >= 2) { return "hi"} else { return "bye"} %></p>`
	s, err := Render(input, NewContext())
	r.NoError(err)
	r.Equal("<p>bye</p>", s)
}

func Test_Render_If_Else_True(t *testing.T) {
	r := require.New(t)
	input := `<p><%= if (true) { %>hi<% } else { %>bye<% } %></p>`
	s, err := Render(input, NewContext())
	r.NoError(err)
	r.Equal("<p>hi</p>", s)
}
