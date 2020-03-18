package mail

import (
	"text/template"

	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/genny/v2/gogen"
	"github.com/gobuffalo/packr/v2"
)

// New mailer generator. It will init the mailers directory if it doesn't already exist
func New(opts *Options) (*genny.Group, error) {
	gg := &genny.Group{}

	if err := opts.Validate(); err != nil {
		return gg, err
	}

	if !opts.SkipInit {
		g, err := initGenerator(opts)
		if err != nil {
			return gg, err
		}
		gg.Add(g)
	}

	g := genny.New()
	h := template.FuncMap{}
	data := map[string]interface{}{
		"opts": opts,
	}
	t := gogen.TemplateTransformer(data, h)
	g.Transformer(t)

	fn := opts.Name.File().String()
	g.File(genny.NewFileS("mailers/"+fn+".go.tmpl", mailerTmpl))
	g.File(genny.NewFileS("templates/mail/"+fn+".plush.html.tmpl", mailTmpl))
	gg.Add(g)

	return gg, nil
}

func initGenerator(opts *Options) (*genny.Generator, error) {
	g := genny.New()

	g.Box(packr.New("github.com/gobuffalo/buffalo/genny/mail/init/templates", "../mail/init/templates"))
	h := template.FuncMap{}
	data := map[string]interface{}{
		"opts": opts,
	}
	t := gogen.TemplateTransformer(data, h)
	g.Transformer(t)

	g.Should = func(r *genny.Runner) bool {
		_, err := r.FindFile("mailers/mailers.go")
		return err != nil
	}
	opts.Name.Titleize()
	return g, nil
}

const mailerTmpl = `package mailers

import (
	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/buffalo/mail"
)

func Send{{.opts.Name.Resource}}() error {
	m := mail.NewMessage()

	// fill in with your stuff:
	m.Subject = "{{.opts.Name.Titleize}}"
	m.From = ""
	m.To = []string{}
	err := m.AddBody(r.HTML("{{.opts.Name.File}}.html"), render.Data{})
	if err != nil {
		return err
	}
	return smtp.Send(m)
}
`

const mailTmpl = `<h2>{{.opts.Name.Titleize}}</h2>

<h3>../templates/mail/{{.opts.Name.File}}.plush.html</h3>`
