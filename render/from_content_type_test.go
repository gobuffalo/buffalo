package render_test

import (
	"bytes"
	"net/http"
	"strings"
	"testing"

	"github.com/gobuffalo/buffalo/render"
	"github.com/stretchr/testify/require"
)

func Test_FromContentType(t *testing.T) {
	r := require.New(t)

	type ji func(v interface{}, req *http.Request) render.Renderer

	type user struct {
		Name string
	}

	table := []ji{
		render.FromContentType,
		render.New(render.Options{}).FromContentType,
	}

	for _, j := range table {

		// Test fallback on JSON
		req, _ := http.NewRequest("GET", "http://localhost/", nil)
		re := j(map[string]string{"hello": "world"}, req)
		r.Equal("application/json", re.ContentType())
		bb := &bytes.Buffer{}
		err := re.Render(bb, nil)
		r.NoError(err)
		r.Equal(`{"hello":"world"}`, strings.TrimSpace(bb.String()))

		// Test format query argument JSON
		req, _ = http.NewRequest("GET", "http://localhost/?format=json", nil)
		re = j(map[string]string{"hello": "world"}, req)
		r.Equal("application/json", re.ContentType())
		bb = &bytes.Buffer{}
		err = re.Render(bb, nil)
		r.NoError(err)
		r.Equal(`{"hello":"world"}`, strings.TrimSpace(bb.String()))

		// Test Content-Type header JSON
		req, _ = http.NewRequest("GET", "http://localhost/", nil)
		req.Header.Set("Content-Type", "application/json")
		re = j(map[string]string{"hello": "world"}, req)
		r.Equal("application/json", re.ContentType())
		bb = &bytes.Buffer{}
		err = re.Render(bb, nil)
		r.NoError(err)
		r.Equal(`{"hello":"world"}`, strings.TrimSpace(bb.String()))

		// Test format query argument XML
		req, _ = http.NewRequest("GET", "http://localhost/?format=xml", nil)
		re = j(user{Name: "mark"}, req)
		r.Equal("application/xml", re.ContentType())
		bb = &bytes.Buffer{}
		err = re.Render(bb, nil)
		r.NoError(err)
		r.Equal("<user>\n  <Name>mark</Name>\n</user>", strings.TrimSpace(bb.String()))

		// Test Content-Type header XML
		req, _ = http.NewRequest("GET", "http://localhost/", nil)
		req.Header.Set("Content-Type", "application/xml")
		re = j(user{Name: "mark"}, req)
		r.Equal("application/xml", re.ContentType())
		bb = &bytes.Buffer{}
		err = re.Render(bb, nil)
		r.NoError(err)
		r.Equal("<user>\n  <Name>mark</Name>\n</user>", strings.TrimSpace(bb.String()))
	}
}
