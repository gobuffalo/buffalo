package updater

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var replace = map[string]string{
	"github.com/markbates/pop":      "github.com/gobuffalo/pop",
	"github.com/markbates/validate": "github.com/gobuffalo/validate",
	"github.com/satori/go.uuid":     "github.com/gobuffalo/uuid",
}

var ic = ImportConverter{
	Data: replace,
}

var checks = []Check{
	ic.Process,
	WebpackCheck,
	PackageJSONCheck,
	DepEnsure,
	DeprecrationsCheck,
}

func ask(q string) bool {
	fmt.Printf("? %s [y/n]\n", q)

	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')

	text = strings.ToLower(strings.TrimSpace(text))
	return text == "y" || text == "yes"
}
