package builder

var tmpl = `package {{.Name}}

import "github.com/gobuffalo/packr"

// !!! GENERATE FILE !!!
// Do NOT hand edit this file!!
// It is recommended that you do not check into this file into SCM.
// We STRONGLY recommend you delete this file after you have built your
// Go binary. You can use the "packr clean" command to clean up this,
// and any other packr generated files.
func init() {
	{{ range $box := .Boxes -}}
	{{range .Files -}}
		packr.PackJSONBytes("{{$box.Name}}", "{{.Name}}", "{{.Contents}}")
	{{end -}}
	{{end -}}
}
`
