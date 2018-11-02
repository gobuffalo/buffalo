package render

import (
	"bytes"
	"strings"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func Test_Template_Partial(t *testing.T) {
	r := require.New(t)

	err := withHTMLFile("index.html", `<%= partial("foo.html") %>`, func(e *Engine) {
		err := withHTMLFile("_foo.html", "Foo > <%= name %>", func(e *Engine) {

			re := e.Template("foo/bar", "index.html")
			r.Equal("foo/bar", re.ContentType())
			bb := &bytes.Buffer{}
			err := re.Render(bb, Data{"name": "Mark"})
			r.NoError(err)
			r.Equal("Foo > Mark", strings.TrimSpace(bb.String()))

		})
		r.NoError(err)
	})
	r.NoError(err)

}

func Test_Template_Partial_WithoutExtension(t *testing.T) {
	r := require.New(t)

	err := withHTMLFile("index.html", `<%= partial("foo") %>`, func(e *Engine) {
		err := withHTMLFile("_foo.html", "Foo > <%= name %>", func(e *Engine) {

			re := e.Template("text/html; charset=utf-8", "index.html")
			bb := &bytes.Buffer{}
			err := re.Render(bb, Data{"name": "Mark"})
			r.NoError(err)
			r.Equal("Foo > Mark", strings.TrimSpace(bb.String()))

		})
		r.NoError(err)
	})
	r.NoError(err)

}

func Test_Template_Partial_Form(t *testing.T) {
	r := require.New(t)

	const newHTML = `<%= form_for(user, {}) { return partial("form.html") } %>`
	const formHTML = `<%= f.InputTag("Name") %>`
	const result = `<form id="-form" method="POST"><div class="form-group"><label>Name</label><input class=" form-control" id="-Name" name="Name" type="text" value="Mark" /></div></form>`

	u := struct {
		Name string
	}{Name: "Mark"}

	err := withHTMLFile("new.html", newHTML, func(e *Engine) {
		err := withHTMLFile("_form.html", formHTML, func(e *Engine) {

			re := e.Template("", "new.html")
			bb := &bytes.Buffer{}
			err := re.Render(bb, Data{"user": u})
			r.NoError(errors.Cause(err))
			r.Equal(result, strings.TrimSpace(bb.String()))

		})
		r.NoError(err)
	})
	r.NoError(err)

}

func Test_Template_Partial_With_For(t *testing.T) {
	r := require.New(t)

	const forHTML = `<%= for(user) in users { %><%= partial("row") %><% } %>`
	const rowHTML = `Hi <%= user.Name %>, `
	const result = `Hi Mark, Hi Yonghwan,`

	users := []struct {
		Name string
	}{{Name: "Mark"}, {Name: "Yonghwan"}}

	err := withHTMLFile("for.html", forHTML, func(e *Engine) {
		err := withHTMLFile("_row.html", rowHTML, func(e *Engine) {

			re := e.Template("text/html; charset=utf-8", "for.html")
			bb := &bytes.Buffer{}
			err := re.Render(bb, Data{"users": users})
			r.NoError(err)
			r.Equal(result, strings.TrimSpace(bb.String()))

		})
		r.NoError(err)
	})
	r.NoError(err)

}

func Test_Template_Partial_With_For_And_Local(t *testing.T) {
	r := require.New(t)

	const forHTML = `<%= for(user) in users { %><%= partial("row", {say:"Hi"}) %><% } %>`
	const rowHTML = `<%= say %> <%= user.Name %>, `
	const result = `Hi Mark, Hi Yonghwan,`

	users := []struct {
		Name string
	}{{Name: "Mark"}, {Name: "Yonghwan"}}

	err := withHTMLFile("for.html", forHTML, func(e *Engine) {
		err := withHTMLFile("_row.html", rowHTML, func(e *Engine) {

			re := e.Template("text/html; charset=utf-8", "for.html")
			bb := &bytes.Buffer{}
			err := re.Render(bb, Data{"users": users})
			r.NoError(err)
			r.Equal(result, strings.TrimSpace(bb.String()))

		})
		r.NoError(err)
	})
	r.NoError(err)

}

func Test_Template_Partial_Recursive_With_Global_And_Local_Context(t *testing.T) {
	r := require.New(t)

	const indexHTML = `<%= partial("foo.html", {other: "Other"}) %>`
	const fooHTML = `<%= other %>|<%= name %>`
	const result = `Other|Mark`

	err := withHTMLFile("index.html", indexHTML, func(e *Engine) {
		err := withHTMLFile("_foo.html", fooHTML, func(e *Engine) {
			re := e.Template("", "index.html")
			bb := &bytes.Buffer{}
			err := re.Render(bb, Data{"name": "Mark"})
			r.NoError(errors.Cause(err))
			r.Equal(result, strings.TrimSpace(bb.String()))
		})
		r.NoError(err)
	})
	r.NoError(err)
}

func Test_Template_Partial_With_Layout(t *testing.T) {
	r := require.New(t)

	err := withHTMLFile("index.html", `<%= partial("foo.html",{layout:"layout.html"}) %>`, func(e *Engine) {
		err := withHTMLFile("_layout.html", "Layout > <%= yield %>", func(e *Engine) {
			err := withHTMLFile("_foo.html", "Foo > <%= name %>", func(e *Engine) {

				re := e.Template("foo/bar", "index.html")
				//r.Equal("foo/bar", re.ContentType())
				bb := &bytes.Buffer{}
				err := re.Render(bb, Data{"name": "Mark"})
				r.NoError(err)
				r.Equal("Layout > Foo > Mark", strings.TrimSpace(bb.String()))

			})
			r.NoError(err)
		})
		r.NoError(err)

	})
	r.NoError(err)

}
