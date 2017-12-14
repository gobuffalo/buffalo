package mail

import (
	"os"
	"path/filepath"

	"github.com/gobuffalo/buffalo/generators"
	"github.com/gobuffalo/buffalo/meta"
	"github.com/gobuffalo/makr"
	"github.com/pkg/errors"
)

// Generator for creating new mailers
type Generator struct {
	App      meta.App  `json:"app"`
	Name     meta.Name `json:"name"`
	SkipInit bool      `json:"skip_init"`
}

// Run the new mailer generator. It will init the mailers directory
// if it doesn't already exist
func (d Generator) Run(root string, data makr.Data) error {
	g := makr.New()
	defer g.Fmt(root)
	data["opts"] = d

	if err := d.initGenerator(data); err != nil {
		return errors.WithStack(err)
	}

	fn := d.Name.File()
	g.Add(makr.NewFile(filepath.Join("mailers", fn+".go"), mailerTmpl))
	g.Add(makr.NewFile(filepath.Join("templates", "mail", fn+".html"), mailTmpl))
	return g.Run(root, data)
}

func (d Generator) initGenerator(data makr.Data) error {
	files, err := generators.Find(filepath.Join(generators.TemplatesPath, "mail", "init"))
	if err != nil {
		return errors.WithStack(err)
	}
	g := makr.New()
	for _, f := range files {
		g.Add(makr.NewFile(f.WritePath, f.Body))
	}

	g.Should = func(data makr.Data) bool {
		if d.SkipInit {
			return false
		}
		if _, err := os.Stat(filepath.Join("mailers", "mailers.go")); err == nil {
			return false
		}
		return true
	}
	return g.Run(".", data)
}

const mailerTmpl = `package mailers

import (
	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/buffalo/mail"
	"github.com/pkg/errors"
)

func Send{{.opts.Name.Model}}() error {
	m := mail.NewMessage()

	// fill in with your stuff:
	m.Subject = "{{.opts.Name.Title}}"
	m.From = ""
	m.To = []string{}
	err := m.AddBody(r.HTML("{{.opts.Name.File}}.html"), render.Data{})
	if err != nil {
		return errors.WithStack(err)
	}
	return smtp.Send(m)
}
`

const mailTmpl = `<h2>{{.opts.Name.Title}}</h2>

<h3>../templates/mail/{{.opts.Name.File}}.html</h3>`
