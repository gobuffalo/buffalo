package form

import (
	"database/sql/driver"
	"fmt"
	"reflect"
	"sync"

	"github.com/gobuffalo/tags"
	"github.com/markbates/inflect"
	"github.com/markbates/validate"
)

type FormFor struct {
	*Form
	Model      interface{}
	name       string
	dashedName string
	reflection reflect.Value
	Errors     *validate.Errors
}

func NewFormFor(model interface{}, opts tags.Options) *FormFor {
	rv := reflect.ValueOf(model)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	name := rv.Type().Name()
	dashedName := inflect.Dasherize(name)

	if opts["id"] == nil {
		opts["id"] = fmt.Sprintf("%s-form", dashedName)
	}

	errors := loadErrors(opts)
	delete(opts, "errors")

	return &FormFor{
		Form:       New(opts),
		Model:      model,
		name:       name,
		dashedName: dashedName,
		reflection: rv,
		Errors:     errors,
	}
}

func loadErrors(opts tags.Options) *validate.Errors {
	errors := validate.NewErrors()
	if opts["errors"] != nil {
		switch t := opts["errors"].(type) {
		default:
			fmt.Printf("Unexpected errors type %T, please\n", t) // %T prints whatever type t has
		case map[string][]string:
			errors = &validate.Errors{
				Errors: opts["errors"].(map[string][]string),
				Lock:   new(sync.RWMutex),
			}
		case *validate.Errors:
			errors = opts["errors"].(*validate.Errors)
		}
	}

	return errors
}

func (f FormFor) CheckboxTag(field string, opts tags.Options) *tags.Tag {
	f.buildOptions(field, opts)
	return f.Form.CheckboxTag(opts)
}

func (f FormFor) InputTag(field string, opts tags.Options) *tags.Tag {
	f.buildOptions(field, opts)
	f.addFormatTag(field, opts)

	return f.Form.InputTag(opts)
}

func (f FormFor) addFormatTag(field string, opts tags.Options) {
	if opts["format"] != nil {
		return
	}

	toff := reflect.TypeOf(f.Model)
	if toff.Kind() == reflect.Ptr {
		toff = toff.Elem()
	}

	if toff.Kind() == reflect.Struct {
		fi, found := toff.FieldByName(field)

		if !found {
			return
		}

		if format, ok := fi.Tag.Lookup("format"); ok && format != "" {
			opts["format"] = format
		}
	}
}

func (f FormFor) RadioButton(field string, opts tags.Options) *tags.Tag {
	f.buildOptions(field, opts)
	return f.Form.RadioButton(opts)
}

func (f FormFor) SelectTag(field string, opts tags.Options) *SelectTag {
	f.buildOptions(field, opts)
	return f.Form.SelectTag(opts)
}

func (f FormFor) TextArea(field string, opts tags.Options) *tags.Tag {
	f.buildOptions(field, opts)
	return f.Form.TextArea(opts)
}

//SubmitTag adds a submit button to the form
func (f FormFor) SubmitTag(value string, opts tags.Options) *tags.Tag {
	return f.Form.SubmitTag(value, opts)
}

func (f FormFor) buildOptions(field string, opts tags.Options) {

	if opts["value"] == nil {
		opts["value"] = f.value(field)
	}

	if opts["name"] == nil {
		opts["name"] = f.findFieldNameFor(field)
	}

	if opts["id"] == nil {
		opts["id"] = fmt.Sprintf("%s-%s", f.dashedName, field)
	}
}

type interfacer interface {
	Interface() interface{}
}

func (f FormFor) value(field string) interface{} {
	fn := f.reflection.FieldByName(field)

	if fn.IsValid() == false {
		return ""
	}

	i := fn.Interface()
	switch t := i.(type) {
	case driver.Valuer:
		value, _ := t.Value()

		if value == nil {
			return ""
		}

		return fmt.Sprintf("%v", value)
	case interfacer:
		return fmt.Sprintf("%v", t.Interface())
	}
	return i
}

func (f FormFor) findFieldNameFor(field string) string {
	ty := reflect.TypeOf(f.Model)

	if ty.Kind() == reflect.Ptr {
		ty = ty.Elem()
	}

	rf, _ := ty.FieldByName(field)

	formDefined := string(rf.Tag.Get("form"))
	if formDefined != "" && formDefined != "-" {
		return formDefined
	}

	schemaDefined := string(rf.Tag.Get("schema"))
	if schemaDefined != "" && schemaDefined != "-" {
		return schemaDefined
	}

	return field
}
