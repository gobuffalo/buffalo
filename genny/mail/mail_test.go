package mail

import (
	"testing"

	"github.com/gobuffalo/flect/name"
	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/genny/gentest"
	"github.com/stretchr/testify/require"
)

func Test_New_NoMailers(t *testing.T) {
	r := require.New(t)
	gg, err := New(&Options{Name: name.New("foo")})
	r.NoError(err)

	run := gentest.NewRunner()
	gg.With(run)
	r.NoError(run.Run())

	res := run.Results()
	r.Len(res.Commands, 0)
	r.Len(res.Files, 4)

	f := res.Files[0]
	r.Equal("mailers/foo.go", f.Name())
	body := f.String()
	r.Contains(body, `err := m.AddBody(r.HTML("foo.html"), render.Data{})`)

	f = res.Files[1]
	r.Equal("mailers/mailers.go", f.Name())

	f = res.Files[2]
	r.Equal("templates/mail/foo.plush.html", f.Name())
	body = f.String()
	r.Contains(body, `<h3>../templates/mail/foo.plush.html</h3>`)

	f = res.Files[3]
	r.Equal("templates/mail/layout.plush.html", f.Name())
}

func Test_New_WithMailers(t *testing.T) {
	r := require.New(t)

	run := gentest.NewRunner()
	run.Disk.Add(genny.NewFileS("mailers/mailers.go", ""))

	gg, err := New(&Options{Name: name.New("foo")})
	r.NoError(err)
	gg.With(run)
	r.NoError(run.Run())

	res := run.Results()
	r.Len(res.Commands, 0)
	r.Len(res.Files, 3)

	f := res.Files[0]
	r.Equal("mailers/foo.go", f.Name())
	body := f.String()
	r.Contains(body, `err := m.AddBody(r.HTML("foo.html"), render.Data{})`)

	f = res.Files[2]
	r.Equal("templates/mail/foo.plush.html", f.Name())
	body = f.String()
	r.Contains(body, `<h3>../templates/mail/foo.plush.html</h3>`)
}
