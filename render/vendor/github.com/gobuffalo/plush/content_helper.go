package plush

import (
	"html/template"

	"github.com/pkg/errors"
)

// ContentFor stores a block of templating code to be re-used later in the template.
/*
	<%= contentFor("buttons") { %>
		<button>hi</button>
	<% } %>
*/
func contentForHelper(name string, help HelperContext) (template.HTML, error) {
	body, err := help.Block()
	if err != nil {
		return "", errors.WithStack(err)
	}
	b := template.HTML(body)
	help.Set(name, b)
	return b, nil
}

// ContentOf retrieves a stored block for templating and renders it.
/*
	<%= contentOf("buttons") %>
*/
func contentOfHelper(name string, help HelperContext) template.HTML {
	if s := help.Value(name); s != nil {
		return s.(template.HTML)
	}
	return ""
}
