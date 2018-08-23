package mail

import (
	"fmt"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/genny/movinglater/gotools"
	"github.com/gobuffalo/packr"
	"github.com/pkg/errors"
)

// New mailer generator. It will init the mailers directory if it doesn't already exist
func New(opts *Options) (*genny.Group, error) {
	gg := &genny.Group{}
	if len(opts.Name.String()) == 0 {
		return gg, errors.New("you must supply a name for your mailer")
	}

	if !opts.SkipInit {
		g, err := initGenerator(opts)
		if err != nil {
			return gg, errors.WithStack(err)
		}
		gg.Add(g)
	}

	g := genny.New()
	h := template.FuncMap{}
	data := map[string]interface{}{
		"opts": opts,
	}
	t := gotools.TemplateTransformer(data, h)
	g.Transformer(t)
	fn := opts.Name.File()
	g.File(genny.NewFile(filepath.Join("mailers", fn+".go.tmpl"), strings.NewReader(mailerTmpl)))
	g.File(genny.NewFile(filepath.Join("templates", "mail", fn+".html.tmpl"), strings.NewReader(mailTmpl)))
	gg.Add(g)

	return gg, nil
}

func initGenerator(opts *Options) (*genny.Generator, error) {
	g := genny.New()

	h := template.FuncMap{}
	data := map[string]interface{}{
		"opts": opts,
	}
	t := gotools.TemplateTransformer(data, h)

	g.RunFn(func(r *genny.Runner) error {
		path := filepath.Join("mailers", "mailers.go")
		_, err := r.FindFile(path)
		fmt.Println("### err ->", err)
		if err == nil {
			return nil
		}
		box := packr.NewBox("../mail/init/templates")
		box.Walk(func(path string, bf packr.File) error {
			f := genny.NewFile(path, bf)
			f, err := t.Transform(f)
			if err != nil {
				return errors.WithStack(err)
			}
			if err := r.File(f); err != nil {
				return errors.WithStack(err)
			}
			return nil
		})
		return nil
	})
	return g, nil
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
