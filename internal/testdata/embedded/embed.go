package embedded

import (
	"embed"
)

//go:embed *
var files embed.FS

func FS() embed.FS {
	return files
}
