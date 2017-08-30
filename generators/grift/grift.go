package grift

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/gobuffalo/buffalo/generators"
	"github.com/gobuffalo/makr"
)

//New allows to create a new grift task generator
func New(data makr.Data) (*makr.Generator, error) {
	g := makr.New()

	files, err := generators.Find(filepath.Join(generators.TemplatesPath, "grift"))
	if err != nil {
		return nil, err
	}

	path := filepath.Join("grifts", data["filename"].(string))
	file := files[0]

	if _, err := os.Stat(path); os.IsNotExist(err) {
		g.Add(makr.NewFile(path, file.Body))
	} else {
		template, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}

		g.Add(makr.NewFile(path, string(template)+existsTmpl))
	}

	return g, nil
}

var existsTmpl = `
{{ if .plainTask -}}
    var _ = Desc("{{.taskName}}", "TODO")
    var _ = Add("{{.taskName}}", func(c *Context) error {
        return nil
    })
{{ else }}
    {{ $last := .last }}
    {{ range $index, $element := .parts }}
        {{ $isLast := eq $index $last }}
        {{ if not $isLast }}
            {{if eq $index 0}}
                var _ = Namespace("{{$element}}", func(){
            {{ else }}
                Namespace("{{$element}}", func(){
            {{end}}
        {{ else }}
            Desc("{{$element}}", "TODO")
            Add("{{$element}}", func(c *Context) error{
                return nil
            })
        {{ end }}
    {{ end }}

    {{ range $index, $element := .parts }}
        {{ if $index }} }) {{ end }}
    {{ end }}

{{ end }}`
