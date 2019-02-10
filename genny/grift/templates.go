package grift

const tmplHeader = `package grifts

import (
	. "github.com/markbates/grift/grift"
)
`

const tmplBody = `
{{ if .opts.Namespaced }}
		{{ range $index, $element := .opts.Parts }}
			{{ if $.opts.Last $element}}
				Desc("{{$element.File}}", "Task Description")
				Add("{{$element.File}}", func(c *Context) error{
						return nil
				})
			{{ else }}
				{{if eq $index 0}}
						var _ = Namespace("{{$element.File}}", func(){
				{{ else }}
						Namespace("{{$element.File}}", func(){
				{{end}}
			{{ end }}
		{{ end }}
		{{ range $index, $element := .opts.Parts }}
				{{ if $index }} }) {{ end }}
		{{ end }}
{{ else }}
		var _ = Desc("{{.opts.Name.File}}", "Task Description")
		var _ = Add("{{.opts.Name.File}}", func(c *Context) error {
				return nil
		})
{{ end }}`
