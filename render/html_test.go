package render_test

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/gobuffalo/buffalo/render"
	"github.com/stretchr/testify/require"
)

func Test_HTML(t *testing.T) {
	r := require.New(t)

	tmpFile, err := ioutil.TempFile("", "test")
	r.NoError(err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.Write([]byte("{{name}}"))
	r.NoError(err)

	type ji func(...string) render.Renderer
	t.Run("without a layout", func(st *testing.T) {
		r := require.New(st)

		table := []ji{
			render.HTML,
			render.New(render.Options{}).HTML,
		}

		for _, j := range table {
			re := j(tmpFile.Name())
			r.Equal("text/html", re.ContentType())
			bb := &bytes.Buffer{}
			err = re.Render(bb, map[string]interface{}{"name": "Mark"})
			r.NoError(err)
			r.Equal("Mark", strings.TrimSpace(bb.String()))
		}
	})

	t.Run("with a layout", func(st *testing.T) {
		r := require.New(st)

		layout, err := ioutil.TempFile("", "test")
		r.NoError(err)
		defer os.Remove(layout.Name())

		_, err = layout.Write([]byte("<body>{{yield}}</body>"))
		r.NoError(err)

		re := render.New(render.Options{HTMLLayout: layout.Name()})

		st.Run("using just the HTMLLayout", func(sst *testing.T) {
			r := require.New(sst)
			h := re.HTML(tmpFile.Name())

			r.Equal("text/html", h.ContentType())
			bb := &bytes.Buffer{}
			err = h.Render(bb, map[string]interface{}{"name": "Mark"})
			r.NoError(err)
			r.Equal("<body>Mark</body>", strings.TrimSpace(bb.String()))
		})

		st.Run("overriding the HTMLLayout", func(sst *testing.T) {
			r := require.New(sst)
			nlayout, err := ioutil.TempFile("", "test-layout2")
			r.NoError(err)
			defer os.Remove(nlayout.Name())

			_, err = nlayout.Write([]byte("<html>{{yield}}</html>"))
			r.NoError(err)
			h := re.HTML(tmpFile.Name(), nlayout.Name())

			r.Equal("text/html", h.ContentType())
			bb := &bytes.Buffer{}
			err = h.Render(bb, map[string]interface{}{"name": "Mark"})
			r.NoError(err)
			r.Equal("<html>Mark</html>", strings.TrimSpace(bb.String()))
		})

	})

	t.Run("with an alternative layout", func(st *testing.T) {
		r := require.New(st)

		// first create the default layout
		layout, err := ioutil.TempFile("", "test")
		r.NoError(err)
		defer os.Remove(layout.Name())

		_, err = layout.Write([]byte("<body>{{yield}}</body>"))
		r.NoError(err)

		// first create the alternative layout
		altLayout, err := ioutil.TempFile("", "test")
		r.NoError(err)
		defer os.Remove(altLayout.Name())

		_, err = altLayout.Write([]byte("<body><!--altLayout-->{{yield}}</body>"))
		r.NoError(err)

		re := render.New(render.Options{HTMLLayout: layout.Name()}).HTML(tmpFile.Name(), altLayout.Name())

		r.Equal("text/html", re.ContentType())
		bb := &bytes.Buffer{}
		err = re.Render(bb, map[string]interface{}{"name": "Mark"})
		r.NoError(err)
		r.Equal("<body><!--altLayout-->Mark</body>", strings.TrimSpace(bb.String()))
	})
}
