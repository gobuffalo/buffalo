package render_test

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/packr"
	"github.com/markbates/willie"
	"github.com/stretchr/testify/require"
)

func Test_Auto_JSON(t *testing.T) {
	r := require.New(t)

	ctx := context.WithValue(context.Background(), "contentType", "application/json")

	ir := render.Auto(ctx, []string{"John", "Paul", "George", "Ringo"})
	r.Equal("application/json", ir.ContentType())

	bb := &bytes.Buffer{}
	r.NoError(ir.Render(bb, render.Data{}))
	r.Equal(`["John","Paul","George","Ringo"]`, strings.TrimSpace(bb.String()))
}

func Test_Auto_XML(t *testing.T) {
	r := require.New(t)

	ctx := context.WithValue(context.Background(), "contentType", "application/xml")

	ir := render.Auto(ctx, []string{"John", "Paul", "George", "Ringo"})
	r.Equal("application/xml", ir.ContentType())

	bb := &bytes.Buffer{}
	r.NoError(ir.Render(bb, render.Data{}))
	r.Equal("<string>John</string>\n<string>Paul</string>\n<string>George</string>\n<string>Ringo</string>", strings.TrimSpace(bb.String()))
}

type Beatle struct {
	ID   int
	Name string
}

type Beatles []Beatle

func Test_Auto_HTML_List(t *testing.T) {
	r := require.New(t)
	ctx := context.WithValue(context.Background(), "method", "GET")

	err := withHTMLFile("beatles/index.html", "INDEX: <%= len(beatles) %>", func(e *render.Engine) {
		ir := e.Auto(ctx, Beatles{
			{Name: "John"},
			{Name: "Paul"},
		})
		r.Equal("text/html", ir.ContentType())
		bb := &bytes.Buffer{}
		r.NoError(ir.Render(bb, render.Data{
			"method": "GET",
		}))
		r.Contains(bb.String(), "INDEX: 2")
	})
	r.NoError(err)
}

func Test_Auto_HTML_Show(t *testing.T) {
	r := require.New(t)
	ctx := context.WithValue(context.Background(), "method", "GET")

	err := withHTMLFile("beatles/show.html", "Show: <%= beatle.Name %>", func(e *render.Engine) {
		ir := e.Auto(ctx, Beatle{
			Name: "John",
		})
		r.Equal("text/html", ir.ContentType())
		bb := &bytes.Buffer{}
		r.NoError(ir.Render(bb, render.Data{
			"method":       "GET",
			"current_path": "/beatles/1",
		}))
		r.Contains(bb.String(), "Show: John")
	})
	r.NoError(err)
}

func Test_Auto_HTML_New(t *testing.T) {
	r := require.New(t)
	ctx := context.WithValue(context.Background(), "method", "GET")

	err := withHTMLFile("beatles/new.html", "New: <%= beatle.Name %>", func(e *render.Engine) {
		ir := e.Auto(ctx, Beatle{
			Name: "John",
		})
		r.Equal("text/html", ir.ContentType())
		bb := &bytes.Buffer{}
		r.NoError(ir.Render(bb, render.Data{
			"method":       "GET",
			"current_path": "/beatles/new",
		}))
		r.Contains(bb.String(), "New: John")
	})
	r.NoError(err)
}

func Test_Auto_HTML_Create(t *testing.T) {
	r := require.New(t)
	ctx := context.WithValue(context.Background(), "method", "POST")

	err := withHTMLFile("beatles/new.html", "New: <%= beatle.Name %>", func(e *render.Engine) {
		ir := e.Auto(ctx, Beatle{
			Name: "John",
		})
		r.Equal("text/html", ir.ContentType())
		bb := &bytes.Buffer{}
		r.NoError(ir.Render(bb, render.Data{
			"method":       "POST",
			"current_path": "/beatles",
		}))
		r.Contains(bb.String(), "New: John")
	})
	r.NoError(err)
}

func Test_Auto_HTML_Create_Redirect(t *testing.T) {
	r := require.New(t)

	app := buffalo.New(buffalo.Options{})
	app.POST("/beatles", func(c buffalo.Context) error {
		b := Beatle{
			ID:   1,
			Name: "John",
		}
		return c.Render(302, render.Auto(c, b))
	})

	w := willie.New(app)
	res := w.HTML("/beatles").Post(nil)
	r.Equal("/beatles/1", res.Location())
}

func Test_Auto_HTML_Edit(t *testing.T) {
	r := require.New(t)
	ctx := context.WithValue(context.Background(), "method", "GET")

	err := withHTMLFile("beatles/edit.html", "Edit: <%= beatle.Name %>", func(e *render.Engine) {
		ir := e.Auto(ctx, Beatle{
			Name: "John",
		})
		r.Equal("text/html", ir.ContentType())
		bb := &bytes.Buffer{}
		r.NoError(ir.Render(bb, render.Data{
			"method":       "GET",
			"current_path": "/beatles/1/edit",
		}))
		r.Contains(bb.String(), "Edit: John")
	})
	r.NoError(err)
}

func Test_Auto_HTML_Update(t *testing.T) {
	r := require.New(t)
	ctx := context.WithValue(context.Background(), "method", "PUT")

	err := withHTMLFile("beatles/edit.html", "Update: <%= beatle.Name %>", func(e *render.Engine) {
		ir := e.Auto(ctx, Beatle{
			Name: "John",
		})
		r.Equal("text/html", ir.ContentType())
		bb := &bytes.Buffer{}
		r.NoError(ir.Render(bb, render.Data{
			"method": "PUT",
		}))
		r.Contains(bb.String(), "Update: John")
	})
	r.NoError(err)
}

func Test_Auto_HTML_Update_Redirect(t *testing.T) {
	r := require.New(t)

	app := buffalo.New(buffalo.Options{})
	app.PUT("/beatles/{id}", func(c buffalo.Context) error {
		b := Beatle{
			ID:   1,
			Name: "John",
		}
		return c.Render(302, render.Auto(c, b))
	})

	w := willie.New(app)
	res := w.HTML("/beatles/1").Put(nil)
	r.Equal("/beatles/1", res.Location())
}

func Test_Auto_HTML_Destroy_Redirect(t *testing.T) {
	r := require.New(t)

	app := buffalo.New(buffalo.Options{})
	app.DELETE("/beatles/{id}", func(c buffalo.Context) error {
		b := Beatle{
			ID:   1,
			Name: "John",
		}
		return c.Render(302, render.Auto(c, b))
	})

	w := willie.New(app)
	res := w.HTML("/beatles/1").Delete()
	r.Equal("/beatles", res.Location())
}

func withHTMLFile(name string, contents string, fn func(*render.Engine)) error {
	tmpDir := filepath.Join(os.TempDir(), filepath.Dir(name))
	err := os.MkdirAll(tmpDir, 0766)
	if err != nil {
		return err
	}
	defer os.Remove(tmpDir)

	tmpFile, err := os.Create(filepath.Join(tmpDir, filepath.Base(name)))
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.Write([]byte(contents))
	if err != nil {
		return err
	}

	e := render.New(render.Options{
		TemplatesBox: packr.NewBox(os.TempDir()),
	})

	fn(e)
	return nil
}
