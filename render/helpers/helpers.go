package helpers

import (
	"encoding/json"
	"html/template"

	"github.com/aymerick/raymond"
)

var Helpers = map[string]interface{}{
	"js_escape":   template.JSEscapeString,
	"html_escape": template.HTMLEscapeString,
	"json":        ToJSON,
	"content_for": ContentFor,
	"content_of":  ContentOf,
}

func ToJSON(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		return err.Error()
	}
	return string(b)
}

func ContentFor(name string, options *raymond.Options) string {
	ctx := options.Ctx().(map[string]interface{})
	body := options.Fn()
	ctx[name] = raymond.SafeString(body)
	return ""
}

func ContentOf(name string, options *raymond.Options) raymond.SafeString {
	ctx := options.Ctx().(map[string]interface{})
	if s, ok := ctx[name]; ok {
		return s.(raymond.SafeString)
	}
	return raymond.SafeString("")
}
