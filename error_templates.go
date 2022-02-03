package buffalo

import (
	_ "embed"
)

var (
	//go:embed error.dev.html
	devErrorTmpl string

	//go:embed error.prod.html
	prodErrorTmpl string

	//go:embed notfound.prod.html
	prodNotFoundTmpl string
)
