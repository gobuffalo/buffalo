package render

import (
	"bytes"
	"strings"
	"testing"

	"github.com/gobuffalo/buffalo/internal/errx"
	"github.com/gobuffalo/packd"
	"github.com/stretchr/testify/require"
)

func Test_Template_Partial(t *testing.T) {
	r := require.New(t)

	const indexHTML = `<%= partial("foo.html") %>`
	const fooHTML = "Foo > <%= name %>"

	box := packd.NewMemoryBox()
	r.NoError(box.AddString("index.html", indexHTML))
	r.NoError(box.AddString("_foo.html", fooHTML))

	re := New(Options{
		TemplatesBox: box,
	})

	bb := &bytes.Buffer{}

	err := re.Template("foo/bar", "index.html").Render(bb, Data{"name": "Mark"})
	r.NoError(err)
	r.Equal("Foo > Mark", strings.TrimSpace(bb.String()))

}

func Test_Template_Partial_WithoutExtension(t *testing.T) {
	r := require.New(t)

	const indexHTML = `<%= partial("foo") %>`
	const fooHTML = "Foo > <%= name %>"

	box := packd.NewMemoryBox()
	r.NoError(box.AddString("index.html", indexHTML))
	r.NoError(box.AddString("_foo.html", fooHTML))

	re := New(Options{
		TemplatesBox: box,
	})

	bb := &bytes.Buffer{}

	err := re.Template("text/html", "index.html").Render(bb, Data{"name": "Mark"})
	r.NoError(err)
	r.Equal("Foo > Mark", strings.TrimSpace(bb.String()))
}

func Test_Template_Partial_Form(t *testing.T) {
	r := require.New(t)

	const newHTML = `<%= form_for(user, {}) { return partial("form.html") } %>`
	const formHTML = `<%= f.InputTag("Name") %>`
	const result = `<form action="/Mark" id="widget-form" method="POST"><div class="form-group"><label>Name</label><input class=" form-control" id="widget-Name" name="Name" type="text" value="Mark" /></div></form>`

	u := Widget{Name: "Mark"}

	re := New(Options{
		TemplatesBox: packd.NewMemoryBox(),
	})
	err := re.TemplatesBox.AddString("new.html", newHTML)
	r.NoError(err)

	err = re.TemplatesBox.AddString("_form.html", formHTML)
	r.NoError(err)

	bb := &bytes.Buffer{}
	err = re.HTML("new.html").Render(bb, Data{"user": u})
	r.NoError(errx.Unwrap(err))
	r.Equal(result, strings.TrimSpace(bb.String()))

}

func Test_Template_Partial_With_For(t *testing.T) {
	r := require.New(t)

	const forHTML = `<%= for(user) in users { %><%= partial("row") %><% } %>`
	const rowHTML = `Hi <%= user.Name %>, `
	const result = `Hi Mark, Hi Yonghwan,`

	box := packd.NewMemoryBox()
	r.NoError(box.AddString("for.html", forHTML))
	r.NoError(box.AddString("_row.html", rowHTML))

	re := New(Options{
		TemplatesBox: box,
	})

	bb := &bytes.Buffer{}

	tmpl := re.Template("text/html; charset=utf-8", "for.html")
	r.Equal("text/html; charset=utf-8", tmpl.ContentType())

	err := tmpl.Render(bb, Data{"users": []Widget{
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

	box := packd.NewMemoryBox()
	r.NoError(box.AddString("for.html", forHTML))
	r.NoError(box.AddString("_row.html", rowHTML))

	re := New(Options{
		TemplatesBox: box,
	})

	bb := &bytes.Buffer{}

	tmpl := re.Template("text/html; charset=utf-8", "for.html")
	r.Equal("text/html; charset=utf-8", tmpl.ContentType())

	err := tmpl.Render(bb, Data{"users": []Widget{
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

	box := packd.NewMemoryBox()
	r.NoError(box.AddString("index.html", indexHTML))
	r.NoError(box.AddString("_foo.html", fooHTML))

	re := New(Options{
		TemplatesBox: box,
	})

	bb := &bytes.Buffer{}

	err := re.Template("foo/bar", "index.html").Render(bb, Data{"name": "Mark"})
	r.NoError(err)
	r.Equal(result, strings.TrimSpace(bb.String()))
}

func Test_Template_Partial_With_Layout(t *testing.T) {
	r := require.New(t)

	const indexHTML = `<%= partial("foo.html",{layout:"layout.html"}) %>`
	const layoutHTML = `Layout > <%= yield %>`
	const fooHTML = "Foo > <%= name %>"
	const result = `Layout > Foo > Mark`

	box := packd.NewMemoryBox()
	r.NoError(box.AddString("index.html", indexHTML))
	r.NoError(box.AddString("_layout.html", layoutHTML))
	r.NoError(box.AddString("_foo.html", fooHTML))

	re := New(Options{
		TemplatesBox: box,
	})

	bb := &bytes.Buffer{}

	err := re.Template("foo/bar", "index.html").Render(bb, Data{"name": "Mark"})
	r.NoError(err)
	r.Equal(result, strings.TrimSpace(bb.String()))

}
