package helpers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"strings"

	"github.com/aymerick/raymond"
	"github.com/markbates/inflect"
	"github.com/shurcooL/github_flavored_markdown"
)

// Helpers that are automatically injected into templates.
var Helpers = map[string]interface{}{
	"js_escape":   template.JSEscapeString,
	"html_escape": template.HTMLEscapeString,
	"json":        ToJSON,
	"content_for": ContentFor,
	"content_of":  ContentOf,
	"upcase":      strings.ToUpper,
	"downcase":    strings.ToLower,
	"markdown":    Markdown,
	"debug":       Debug,
}

func init() {
	for k, v := range inflect.Helpers {
		Helpers[k] = v
	}
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

// Markdown converts the string into HTML using GitHub flavored markdown.
func Markdown(body string) raymond.SafeString {
	b := github_flavored_markdown.Markdown([]byte(body))
	return raymond.SafeString(string(b))
}

// Debug by verbosely printing out using 'pre' tags.
func Debug(v interface{}) raymond.SafeString {
	return raymond.SafeString(fmt.Sprintf("<pre>%+v</pre>", v))
}
