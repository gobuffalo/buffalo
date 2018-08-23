package mail

import (
	"context"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gobuffalo/genny"
	"github.com/stretchr/testify/require"
)

func Test_New_NoMailers(t *testing.T) {
	r := require.New(t)
	gg, err := New(&Options{Name: "foo"})
	r.NoError(err)

	run := genny.DryRunner(context.Background())
	gg.With(run)
	r.NoError(run.Run())

	res := run.Results()
	r.Len(res.Commands, 0)
	r.Len(res.Files, 4)

	f := res.Files[0]
	r.Equal(filepath.Join("mailers", "foo.go"), f.Name())
	body := f.String()
	r.Contains(body, `err := m.AddBody(r.HTML("foo.html"), render.Data{})`)

	f = res.Files[1]
	r.Equal(filepath.Join("mailers", "mailers.go"), f.Name())

	f = res.Files[2]
	r.Equal(filepath.Join("templates", "mail", "foo.html"), f.Name())
	body = f.String()
	r.Contains(body, `<h3>../templates/mail/foo.html</h3>`)

	f = res.Files[3]
	r.Equal(filepath.Join("templates", "mail", "layout.html"), f.Name())
}

func Test_New_WithMailers(t *testing.T) {
	r := require.New(t)

	run := genny.DryRunner(context.Background())
	g := genny.New()
	g.File(genny.NewFile(filepath.Join("mailers", "mailers.go"), strings.NewReader("")))
	run.With(g)

	gg, err := New(&Options{Name: "foo"})
	r.NoError(err)
	gg.With(run)
	r.NoError(run.Run())

	res := run.Results()
	r.Len(res.Commands, 0)
	r.Len(res.Files, 3)

	f := res.Files[0]
	r.Equal(filepath.Join("mailers", "foo.go"), f.Name())
	body := f.String()
	r.Contains(body, `err := m.AddBody(r.HTML("foo.html"), render.Data{})`)

	f = res.Files[2]
	r.Equal(filepath.Join("templates", "mail", "foo.html"), f.Name())
	body = f.String()
	r.Contains(body, `<h3>../templates/mail/foo.html</h3>`)

}
