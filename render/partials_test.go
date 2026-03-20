package render

import (
	"bytes"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/require"
)

func Test_Template_Partial(t *testing.T) {
	r := require.New(t)

	const part = `<%= partial("foo.html") %>`
	const tmpl = "Foo > <%= name %>"

	rootFS := fstest.MapFS{
		htmlTemplate: &fstest.MapFile{
			Data: []byte(tmpl),
			Mode: 0644,
		},
		"_foo.html": &fstest.MapFile{
			Data: []byte(part),
			Mode: 0644,
		},
	}

	e := NewEngine()
	e.TemplatesFS = rootFS

	bb := &bytes.Buffer{}

	re := e.Template("foo/bar", htmlTemplate)
	r.NoError(re.Render(bb, Data{"name": "Mark"}))
	r.Equal("Foo > Mark", strings.TrimSpace(bb.String()))
}

func Test_Template_PartialCustomFeeder(t *testing.T) {
	r := require.New(t)

	rootFS := fstest.MapFS{
		"base.plush.html": &fstest.MapFile{
			Data: []byte(`<%= partial("foo.plush.html") %>`),
			Mode: 0644,
		},
		"_foo.plush.html": &fstest.MapFile{
			Data: []byte("other"),
			Mode: 0644,
		},
	}

	e := NewEngine()
	e.TemplatesFS = rootFS

	t.Run("Custom Feeder", func(t *testing.T) {
		e.Helpers["partialFeeder"] = func(path string) (string, error) {
			return "custom", nil
		}

		bb := &bytes.Buffer{}

		re := e.HTML("base.plush.html")
		r.NoError(re.Render(bb, Data{}))
		r.Equal("custom", strings.TrimSpace(bb.String()))
	})

	t.Run("Default Feeder", func(t *testing.T) {
		e.Helpers["partialFeeder"] = nil

		bb := &bytes.Buffer{}

		re := e.HTML("base.plush.html")
		r.NoError(re.Render(bb, Data{}))
		r.Equal("other", strings.TrimSpace(bb.String()))
	})
}

func Test_Template_Partial_WithoutExtension(t *testing.T) {
	r := require.New(t)

	const part = `<%= partial("foo") %>`
	const tmpl = "Foo > <%= name %>"

	rootFS := fstest.MapFS{
		htmlTemplate: &fstest.MapFile{
			Data: []byte(tmpl),
			Mode: 0644,
		},
		"_foo.html": &fstest.MapFile{
			Data: []byte(part),
			Mode: 0644,
		},
	}

	e := NewEngine()
	e.TemplatesFS = rootFS

	bb := &bytes.Buffer{}

	re := e.Template("foo/bar", htmlTemplate)
	r.NoError(re.Render(bb, Data{"name": "Mark"}))
	r.Equal("Foo > Mark", strings.TrimSpace(bb.String()))
}

func Test_Template_Partial_Form(t *testing.T) {
	r := require.New(t)

	const newHTML = `<%= form_for(user, {}) { return partial("form.html") } %>`
	const formHTML = `<%= f.InputTag("Name") %>`
	const result = `<form action="/Mark" id="widget-form" method="POST"><div class="form-group"><label class="form-label" for="widget-Name">Name</label><input class="form-control" id="widget-Name" name="Name" type="text" value="Mark" /></div></form>`

	rootFS := fstest.MapFS{
		"new.html": &fstest.MapFile{
			Data: []byte(newHTML),
			Mode: 0644,
		},
		"_form.html": &fstest.MapFile{
			Data: []byte(formHTML),
			Mode: 0644,
		},
	}

	e := NewEngine()
	e.TemplatesFS = rootFS

	u := Widget{Name: "Mark"}

	bb := &bytes.Buffer{}
	re := e.HTML("new.html")
	r.NoError(re.Render(bb, Data{"user": u}))
	r.Equal(result, strings.TrimSpace(bb.String()))
}

func Test_Template_Partial_With_For(t *testing.T) {
	r := require.New(t)

	const forHTML = `<%= for(user) in users { %><%= partial("row") %><% } %>`
	const rowHTML = `Hi <%= user.Name %>, `
	const result = `Hi Mark, Hi Yonghwan,`

	rootFS := fstest.MapFS{
		"for.html": &fstest.MapFile{
			Data: []byte(forHTML),
			Mode: 0644,
		},
		"_row.html": &fstest.MapFile{
			Data: []byte(rowHTML),
			Mode: 0644,
		},
	}

	e := NewEngine()
	e.TemplatesFS = rootFS

	bb := &bytes.Buffer{}

	re := e.Template("text/html; charset=utf-8", "for.html")
	r.Equal("text/html; charset=utf-8", re.ContentType())

	err := re.Render(bb, Data{"users": []Widget{
		{Name: "Mark"},
		{Name: "Yonghwan"},
	}})

	r.NoError(err)
	r.Equal(result, strings.TrimSpace(bb.String()))
}

func Test_Template_Partial_With_For_And_Local(t *testing.T) {
	r := require.New(t)

	const forHTML = `<%= for(user) in users { %><%= partial("row", {say:"Hi"}) %><% } %>`
	const rowHTML = `<%= say %> <%= user.Name %>, `
	const result = `Hi Mark, Hi Yonghwan,`

	rootFS := fstest.MapFS{
		"for.html": &fstest.MapFile{
			Data: []byte(forHTML),
			Mode: 0644,
		},
		"_row.html": &fstest.MapFile{
			Data: []byte(rowHTML),
			Mode: 0644,
		},
	}

	e := NewEngine()
	e.TemplatesFS = rootFS

	bb := &bytes.Buffer{}

	re := e.Template("text/html; charset=utf-8", "for.html")
	r.Equal("text/html; charset=utf-8", re.ContentType())

	err := re.Render(bb, Data{"users": []Widget{
		{Name: "Mark"},
		{Name: "Yonghwan"},
	}})

	r.NoError(err)
	r.Equal(result, strings.TrimSpace(bb.String()))
}

func Test_Template_Partial_Recursive_With_Global_And_Local_Context(t *testing.T) {
	r := require.New(t)

	const indexHTML = `<%= partial("foo.html", {other: "Other"}) %>`
	const fooHTML = `<%= other %>|<%= name %>`
	const result = `Other|Mark`

	rootFS := fstest.MapFS{
		"index.html": &fstest.MapFile{
			Data: []byte(indexHTML),
			Mode: 0644,
		},
		"_foo.html": &fstest.MapFile{
			Data: []byte(fooHTML),
			Mode: 0644,
		},
	}

	e := NewEngine()
	e.TemplatesFS = rootFS

	bb := &bytes.Buffer{}

	re := e.Template("foo/bar", "index.html")
	r.NoError(re.Render(bb, Data{"name": "Mark"}))
	r.Equal(result, strings.TrimSpace(bb.String()))
}

func Test_Template_Partial_With_Layout(t *testing.T) {
	r := require.New(t)

	const indexHTML = `<%= partial("foo.html",{layout:"layout.html"}) %>`
	const layoutHTML = `Layout > <%= yield %>`
	const fooHTML = "Foo > <%= name %>"
	const result = `Layout > Foo > Mark`

	rootFS := fstest.MapFS{
		"index.html": &fstest.MapFile{
			Data: []byte(indexHTML),
			Mode: 0644,
		},
		"_layout.html": &fstest.MapFile{
			Data: []byte(layoutHTML),
			Mode: 0644,
		},
		"_foo.html": &fstest.MapFile{
			Data: []byte(fooHTML),
			Mode: 0644,
		},
	}

	e := NewEngine()
	e.TemplatesFS = rootFS

	bb := &bytes.Buffer{}

	re := e.Template("foo/bar", "index.html")
	r.NoError(re.Render(bb, Data{"name": "Mark"}))
	r.Equal(result, strings.TrimSpace(bb.String()))
}
