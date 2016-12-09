package helpers

import (
	"testing"

	"github.com/aymerick/raymond"
	"github.com/stretchr/testify/require"
)

func Test_ToJSON(t *testing.T) {
	r := require.New(t)
	s := ToJSON([]string{"mark", "bates"})
	r.Equal(`["mark","bates"]`, s)
}

func Test_ContentForOf(t *testing.T) {
	r := require.New(t)
	html := `
	{{#content_for "buttons"}}<button>hi</button>{{/content_for}}
	<b1>{{content_of "buttons"}}</b1>
	<b2>{{content_of "buttons"}}</b2>
	`
	tmpl := raymond.MustParse(html)
	tmpl.RegisterHelpers(Helpers)

	body := tmpl.MustExec(map[string]interface{}{})
	r.Contains(body, "<b1><button>hi</button></b1>")
	r.Contains(body, "<b2><button>hi</button></b2>")
}
