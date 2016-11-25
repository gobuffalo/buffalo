package render

import "html/template"

type Options struct {
	HTMLLayout    string
	TemplatesPath string
	templates     *template.Template
	TemplateFuncs template.FuncMap
}
