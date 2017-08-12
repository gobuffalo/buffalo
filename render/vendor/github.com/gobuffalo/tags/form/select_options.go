package form

import (
	"bytes"
	"html/template"
)

type SelectOption struct {
	Value         interface{}
	Label         interface{}
	SelectedValue interface{}
}

func (s SelectOption) String() string {
	v := template.HTMLEscaper(s.Value)
	sv := template.HTMLEscaper(s.SelectedValue)
	l := template.HTMLEscaper(s.Label)
	bb := &bytes.Buffer{}
	bb.WriteString(`<option value="`)
	bb.WriteString(v)
	bb.WriteString(`"`)
	if v == sv {
		bb.WriteString(` selected`)
	}
	bb.WriteString(`>`)
	bb.WriteString(l)
	bb.WriteString("</option>")
	return bb.String()
}

type SelectOptions []SelectOption
