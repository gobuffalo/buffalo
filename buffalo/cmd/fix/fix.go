package fix

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

//YesToAll will be used by the command to skip the questions
var YesToAll bool

var replace = map[string]string{
	"github.com/markbates/pop":                     "github.com/gobuffalo/pop",
	"github.com/markbates/validate":                "github.com/gobuffalo/validate",
	"github.com/satori/go.uuid":                    "github.com/gobuffalo/uuid",
	"github.com/markbates/willie":                  "github.com/gobuffalo/httptest",
	"github.com/shurcooL/github_flavored_markdown": "github.com/gobuffalo/github_flavored_markdown",
}

var ic = ImportConverter{
	Data: replace,
}

var mr = MiddlewareTransformer{
	PackagesReplacement: map[string]string{
		"github.com/gobuffalo/buffalo/middleware/basicauth": "github.com/gobuffalo/mw-basicauth",
		"github.com/gobuffalo/buffalo/middleware/csrf":      "github.com/gobuffalo/mw-csrf",
		"github.com/gobuffalo/buffalo/middleware/i18n":      "github.com/gobuffalo/mw-i18n",
		"github.com/gobuffalo/buffalo/middleware/ssl":       "github.com/gobuffalo/mw-forcessl",
		"github.com/gobuffalo/buffalo/middleware/tokenauth": "github.com/gobuffalo/mw-tokenauth",
	},

	Aliases: map[string]string{
		"github.com/gobuffalo/mw-basicauth":          "basicauth",
		"github.com/gobuffalo/mw-csrf":               "csrf",
		"github.com/gobuffalo/mw-i18n":               "i18n",
		"github.com/gobuffalo/mw-forcessl":           "forcessl",
		"github.com/gobuffalo/mw-tokenauth":          "tokenauth",
		"github.com/gobuffalo/mw-paramlogger":        "paramlogger",
		"github.com/gobuffalo/mw-contenttype":        "contenttype",
		"github.com/gobuffalo/buffalo-pop/pop/popmw": "popmw",
	},
}

var checks = []Check{
	PackrClean,
	ic.Process,
	mr.transformPackages,
	WebpackCheck,
	PackageJSONCheck,
	DepEnsure,
	installTools,
	DeprecrationsCheck,
}

func ask(q string) bool {
	if YesToAll {
		return true
	}

	fmt.Printf("? %s [y/n]\n", q)

	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')

	text = strings.ToLower(strings.TrimSpace(text))
	return text == "y" || text == "yes"
}
