package plush

import (
	"html/template"
	"strings"
	"testing"

	"github.com/gobuffalo/tags"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func Test_Render_Simple_HTML(t *testing.T) {
	r := require.New(t)

	input := `<p>Hi</p>`
	s, err := Render(input, NewContext())
	r.NoError(err)
	r.Equal(input, s)
}

func Test_Render_Keeps_Spacing(t *testing.T) {
	r := require.New(t)
	input := `<%= greet %> <%= name %>`

	ctx := NewContext()
	ctx.Set("greet", "hi")
	ctx.Set("name", "mark")

	s, err := Render(input, ctx)
	r.NoError(err)
	r.Equal("hi mark", s)
}

// support identifiers containing digits, but not starting with a digits
func Test_Identifiers_With_Digits(t *testing.T) {
	r := require.New(t)
	input := `<%= my123greet %> <%= name3 %>`

	ctx := NewContext()
	ctx.Set("my123greet", "hi")
	ctx.Set("name3", "mark")

	s, err := Render(input, ctx)
	r.NoError(err)
	r.Equal("hi mark", s)
}

func Test_Render_HTML_InjectedString(t *testing.T) {
	r := require.New(t)

	input := `<p><%= "mark" %></p>`
	s, err := Render(input, NewContext())
	r.NoError(err)
	r.Equal("<p>mark</p>", s)
}

func Test_Render_EscapedString(t *testing.T) {
	r := require.New(t)

	input := `<p><%= "<script>alert('pwned')</script>" %></p>`
	s, err := Render(input, NewContext())
	r.NoError(err)
	r.Equal("<p>&lt;script&gt;alert(&#39;pwned&#39;)&lt;/script&gt;</p>", s)
}

func Test_Render_Injected_Variable(t *testing.T) {
	r := require.New(t)

	input := `<p><%= name %></p>`
	s, err := Render(input, NewContextWith(map[string]interface{}{
		"name": "Mark",
	}))
	r.NoError(err)
	r.Equal("<p>Mark</p>", s)
}

func Test_Render_Let_Hash(t *testing.T) {
	r := require.New(t)

	input := `<p><% let h = {"a": "A"} %><%= h["a"] %></p>`
	s, err := Render(input, NewContext())
	r.NoError(err)
	r.Equal("<p>A</p>", s)
}

func Test_Render_Hash_Array_Index(t *testing.T) {
	r := require.New(t)

	input := `<%= m["first"] + " " + m["last"] %>|<%= a[0+1] %>`
	s, err := Render(input, NewContextWith(map[string]interface{}{
		"m": map[string]string{"first": "Mark", "last": "Bates"},
		"a": []string{"john", "paul"},
	}))
	r.NoError(err)
	r.Equal("Mark Bates|paul", s)
}

func Test_Render_Missing_Variable(t *testing.T) {
	r := require.New(t)

	input := `<p><%= name %></p>`
	s, err := Render(input, NewContext())
	r.NoError(err)
	r.Equal("<p></p>", s)
}

func Test_Render_HTML_Escape(t *testing.T) {
	r := require.New(t)

	input := `<%= escapedHTML() %>|<%= unescapedHTML() %>|<%= raw("<b>unsafe</b>") %>`
	s, err := Render(input, NewContextWith(map[string]interface{}{
		"escapedHTML": func() string {
			return "<b>unsafe</b>"
		},
		"unescapedHTML": func() template.HTML {
			return "<b>unsafe</b>"
		},
	}))
	r.NoError(err)
	r.Equal("&lt;b&gt;unsafe&lt;/b&gt;|<b>unsafe</b>|<b>unsafe</b>", s)
}

func Test_Render_ShowNoShow(t *testing.T) {
	r := require.New(t)
	input := `<%= "shown" %><% "notshown" %>`
	s, err := Render(input, NewContext())
	r.NoError(err)
	r.Equal("shown", s)
}

func Test_Render_Struct_Attribute(t *testing.T) {
	r := require.New(t)
	input := `<%= f.Name %>`
	ctx := NewContext()
	f := struct {
		Name string
	}{"Mark"}
	ctx.Set("f", f)
	s, err := Render(input, ctx)
	r.NoError(err)
	r.Equal("Mark", s)
}

func Test_Render_ScriptFunction(t *testing.T) {
	r := require.New(t)

	input := `<% let add = fn(x) { return x + 2; }; %><%= add(2) %>`

	s, err := Render(input, NewContext())
	if err != nil {
		r.NoError(err)
	}
	r.Equal("4", s)
}

func Test_Render_HasBlock(t *testing.T) {
	r := require.New(t)
	ctx := NewContext()
	ctx.Set("blockCheck", func(help HelperContext) string {
		if help.HasBlock() {
			s, _ := help.Block()
			return s
		}
		return "no block"
	})
	input := `<%= blockCheck() {return "block"} %>|<%= blockCheck() %>`
	s, err := Render(input, ctx)
	r.NoError(err)
	r.Equal("block|no block", s)
}

func Test_Render_HashCall(t *testing.T) {
	r := require.New(t)
	input := `<%= m["a"] %>`
	ctx := NewContext()
	ctx.Set("m", map[string]string{
		"a": "A",
	})
	s, err := Render(input, ctx)
	r.NoError(err)
	r.Equal("A", s)
}

func Test_Render_HashCall_OnAttribute(t *testing.T) {
	r := require.New(t)
	input := `<%= m.MyMap[key] %>`
	ctx := NewContext()
	ctx.Set("m", struct {
		MyMap map[string]string
	}{
		MyMap: map[string]string{"a": "A"},
	})
	ctx.Set("key", "a")
	s, err := Render(input, ctx)
	r.NoError(err)
	r.Equal("A", s)
}

func Test_Render_HashCall_OnAttribute_IntoFunction(t *testing.T) {
	r := require.New(t)
	input := `<%= debug(m.MyMap[key]) %>`
	ctx := NewContext()
	ctx.Set("m", struct {
		MyMap map[string]string
	}{
		MyMap: map[string]string{"a": "A"},
	})
	ctx.Set("key", "a")
	s, err := Render(input, ctx)
	r.NoError(err)
	r.Equal("<pre>A</pre>", s)
}

func Test_Render_UnknownAttribute_on_Callee(t *testing.T) {
	r := require.New(t)
	ctx := NewContext()
	ctx.Set("m", struct{}{})
	input := `<%= m.Foo %>`
	_, err := Render(input, ctx)
	r.Error(err)
	r.Contains(err.Error(), "m.Foo")
}

type Robot struct {
	Avatar Avatar
}

type Avatar string

func (a Avatar) URL() string {
	return strings.ToUpper(string(a))
}

func Test_Render_Function_on_sub_Struct(t *testing.T) {
	r := require.New(t)
	ctx := NewContext()
	bender := Robot{
		Avatar: Avatar("bender.jpg"),
	}
	ctx.Set("robot", bender)
	input := `<%= robot.Avatar.URL() %>`
	s, err := Render(input, ctx)
	r.NoError(err)
	r.Equal("BENDER.JPG", s)
}

func Test_Render_Var_ends_in_Number(t *testing.T) {
	r := require.New(t)
	ctx := NewContextWith(map[string]interface{}{
		"myvar1": []string{"john", "paul"},
	})
	s, err := Render(`<%= for (n) in myvar1 {return n}`, ctx)
	r.NoError(err)
	r.Equal("johnpaul", s)
}

func Test_Render_Dash_in_Helper(t *testing.T) {
	r := require.New(t)
	ctx := NewContextWith(map[string]interface{}{
		"my-helper": func() string {
			return "hello"
		},
	})
	s, err := Render(`<%= my-helper() %>`, ctx)
	r.NoError(err)
	r.Equal("hello", s)
}

func Test_Let_Inside_Helper(t *testing.T) {
	r := require.New(t)
	ctx := NewContextWith(map[string]interface{}{
		"divwrapper": func(opts map[string]interface{}, helper HelperContext) (template.HTML, error) {
			body, err := helper.Block()
			if err != nil {
				return template.HTML(""), errors.WithStack(err)
			}
			t := tags.New("div", opts)
			t.Append(body)
			return t.HTML(), nil
		},
	})

	input := `<%= divwrapper({"class": "myclass"}) { %>
<ul>
    <% let a = [1, 2, "three", "four"] %>
    <%= for (index, name) in a { %>
        <li><%=index%> - <%=name%></li>
    <% } %>
</ul>
<% } %>`

	s, err := Render(input, ctx)
	r.NoError(err)
	r.Contains(s, "<li>0 - 1</li>")
	r.Contains(s, "<li>1 - 2</li>")
	r.Contains(s, "<li>2 - three</li>")
	r.Contains(s, "<li>3 - four</li>")
}
