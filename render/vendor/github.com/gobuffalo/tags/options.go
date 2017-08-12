package tags

import (
	"fmt"
	"html/template"
	"sort"
	"strings"
)

type Options map[string]interface{}

func (o Options) String() string {
	var out = make([]string, 0, len(o))
	var tmp = make([]string, 2)
	for k, v := range o {
		tmp[0] = template.HTMLEscaper(k)
		tmp[1] = fmt.Sprintf("\"%s\"", template.HTMLEscaper(v))
		out = append(out, strings.Join(tmp, "="))
	}
	sort.Strings(out)
	return strings.Join(out, " ")
}
