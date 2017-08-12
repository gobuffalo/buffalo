package plush

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Render_For_Array(t *testing.T) {
	r := require.New(t)
	input := `<% for (i,v) in ["a", "b", "c"] {return v} %>`
	s, err := Render(input, NewContext())
	r.NoError(err)
	r.Equal("", s)
}

func Test_Render_For_Hash(t *testing.T) {
	r := require.New(t)
	input := `<%= for (k,v) in myMap { %><%= k + ":" + v%><% } %>`
	s, err := Render(input, NewContextWith(map[string]interface{}{
		"myMap": map[string]string{
			"a": "A",
			"b": "B",
		},
	}))
	r.NoError(err)
	r.Contains(s, "a:A")
	r.Contains(s, "b:B")
}

func Test_Render_For_Array_Return(t *testing.T) {
	r := require.New(t)
	input := `<%= for (i,v) in ["a", "b", "c"] {return v} %>`
	s, err := Render(input, NewContext())
	r.NoError(err)
	r.Equal("abc", s)
}

func Test_Render_For_Array_Key_Only(t *testing.T) {
	r := require.New(t)
	input := `<%= for (v) in ["a", "b", "c"] {%><%=v%><%} %>`
	s, err := Render(input, NewContext())
	r.NoError(err)
	r.Equal("abc", s)
}

func Test_Render_For_Func_Range(t *testing.T) {
	r := require.New(t)
	input := `<%= for (v) in range(3,5) { %><%=v%><% } %>`
	s, err := Render(input, NewContext())
	r.NoError(err)
	r.Equal("345", s)
}

func Test_Render_For_Func_Between(t *testing.T) {
	r := require.New(t)
	input := `<%= for (v) in between(3,6) { %><%=v%><% } %>`
	s, err := Render(input, NewContext())
	r.NoError(err)
	r.Equal("45", s)
}

func Test_Render_For_Func_Until(t *testing.T) {
	r := require.New(t)
	input := `<%= for (v) in until(3) { %><%=v%><% } %>`
	s, err := Render(input, NewContext())
	r.NoError(err)
	r.Equal("012", s)
}

func Test_Render_For_Array_Key_Value(t *testing.T) {
	r := require.New(t)
	input := `<%= for (i,v) in ["a", "b", "c"] {%><%=i%><%=v%><%} %>`
	s, err := Render(input, NewContext())
	r.NoError(err)
	r.Equal("0a1b2c", s)
}

func Test_Render_For_Nil(t *testing.T) {
	r := require.New(t)
	input := `<% for (i,v) in nilValue {return v} %>`
	ctx := NewContext()
	ctx.Set("nilValue", nil)
	s, err := Render(input, ctx)
	r.NoError(err)
	r.Equal("", s)
}

func Test_Render_For_Map_Nil_Value(t *testing.T) {
	r := require.New(t)
	input := `
	<%= for (k, v) in flash["errors"] { %>
		Flash:
			<%= k %>:<%= v %>
	<% } %>
`
	ctx := NewContext()
	ctx.Set("flash", map[string][]string{})
	s, err := Render(input, ctx)
	r.NoError(err)
	r.Equal("", strings.TrimSpace(s))
}
