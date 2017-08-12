package generate

const modelTemplate = `package {{.package_name}}

import (
	{{ range $i := .model.Imports -}}
	"{{$i}}"
	{{ end -}}
	{{ if .model.ValidatableAttributes -}}
	"github.com/markbates/validate/validators"
	{{ end -}}
)

type {{.model_name}} struct {
	{{range $a := .model.Attributes -}}
	{{$a}}
	{{end -}}
}

// String is not required by pop and may be deleted
func ({{.char}} {{.model_name}}) String() string {
	j{{.char}}, _ := json.Marshal({{.char}})
	return string(j{{.char}})
}

// {{.plural_model_name}} is not required by pop and may be deleted
type {{.plural_model_name}} []{{.model_name}}

// String is not required by pop and may be deleted
func ({{.char}} {{.plural_model_name}}) String() string {
	j{{.char}}, _ := json.Marshal({{.char}})
	return string(j{{.char}})
}

// Validate gets run every time you call a "pop.Validate" method.
// This method is not required and may be deleted.
func ({{.char}} *{{.model_name}}) Validate(tx *pop.Connection) (*validate.Errors, error) {
	{{ if .model.ValidatableAttributes -}}
	return validate.Validate(
		{{ range $a := .model.ValidatableAttributes -}}
		&validators.{{capitalize $a.GoType}}IsPresent{Field: {{$.char}}.{{$a.Names.Proper}}, Name: "{{$a.Names.Proper}}"},
		{{end -}}
	), nil
	{{ else -}}
		return validate.NewErrors(), nil
	{{ end -}}
}

// ValidateSave gets run every time you call "pop.ValidateSave" method.
// This method is not required and may be deleted.
func ({{.char}} *{{.model_name}}) ValidateSave(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateUpdate" method.
// This method is not required and may be deleted.
func ({{.char}} *{{.model_name}}) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
`

const modelTestTemplate = `package {{.package_name}}_test

import "testing"

func Test_{{.model_name}}(t *testing.T) {
	t.Fatal("This test needs to be implemented!")
}
`
