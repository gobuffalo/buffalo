package helpers

import (
	"encoding/json"
	"html/template"
	"strings"

	"github.com/aymerick/raymond"
)

// Helpers that are automatically injected into templates.
/*
yield - renders the content of a template into a layout.
partial - renders the content of a partial into a template.
js_escape - escapes a string to be valid in JavaScript.
html_escape - escapes any HTML characters in a string.
json - converts the interface to JSON.
content_for - stores a block of templating code to be re-used later in the template.
content_of - retrieves a stored block for templating and renders it.
upcase - strings.ToUpper.
downcase - strings.ToLower.
*/
var Helpers = map[string]interface{}{
	"js_escape":   template.JSEscapeString,
	"html_escape": template.HTMLEscapeString,
	"json":        ToJSON,
	"content_for": ContentFor,
	"content_of":  ContentOf,
	"upcase":      strings.ToUpper,
	"downcase":    strings.ToLower,
}

// ToJSON converts an interface into a string.
func ToJSON(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		return err.Error()
	}
	return string(b)
}

// ContentFor stores a block of templating code to be re-used later in the template.
/*
	{{content_for "buttons"}}
		<button>hi</button>
	{{/content_for}}
*/
func ContentFor(name string, options *raymond.Options) string {
	ctx := options.Ctx().(map[string]interface{})
	body := options.Fn()
	ctx[name] = raymond.SafeString(body)
	return ""
}

// ContentOf retrieves a stored block for templating and renders it.
/*
	{{content_of "buttons"}}
*/
func ContentOf(name string, options *raymond.Options) raymond.SafeString {
	ctx := options.Ctx().(map[string]interface{})
	if s, ok := ctx[name]; ok {
		return s.(raymond.SafeString)
	}
	return raymond.SafeString("")
}
