package render

import "html/template"

type Options struct {
	HTMLLayout    string
	templates     *template.Template
	TemplateFuncs template.FuncMap
}
