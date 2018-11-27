package mail

import (
	"text/template"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/genny/movinglater/gotools"
	"github.com/gobuffalo/packr/v2"
	"github.com/pkg/errors"
)

// New mailer generator. It will init the mailers directory if it doesn't already exist
func New(opts *Options) (*genny.Group, error) {
	gg := &genny.Group{}

	if err := opts.Validate(); err != nil {
		return gg, errors.WithStack(err)
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
	g.File(genny.NewFileS("mailers/"+fn+".go.tmpl", mailerTmpl))
	g.File(genny.NewFileS("templates/mail/"+fn+".html.tmpl", mailTmpl))
	gg.Add(g)

	return gg, nil
}

func initGenerator(opts *Options) (*genny.Generator, error) {
	g := genny.New()

	g.Box(packr.New("buffalo:genny:mail:init", "../mail/init/templates"))
	h := template.FuncMap{}
	data := map[string]interface{}{
		"opts": opts,
	}
	t := gotools.TemplateTransformer(data, h)
	g.Transformer(t)

	g.Should = func(r *genny.Runner) bool {
		_, err := r.FindFile("mailers/mailers.go")
		return err != nil
	}
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
