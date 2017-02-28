package render

import (
	"bytes"
	"html/template"
)

// TemplateEngine needs to be implemented for a temlating system be able to be used with Buffalo.
type TemplateEngine func(input string, data map[string]interface{}, helpers map[string]interface{}) (string, error)

// GoTemplateEngine implements the TemplateEngine interface for using standard Go templates
func GoTemplateEngine(input string, data map[string]interface{}, helpers map[string]interface{}) (string, error) {
	t, err := template.New(input).Parse(input)
	if err != nil {
		return "", err
	}
	if helpers != nil {
		t = t.Funcs(helpers)
	}
	bb := &bytes.Buffer{}
	err = t.Execute(bb, data)
	return bb.String(), err
}
