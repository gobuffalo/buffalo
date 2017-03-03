package refresh

import "github.com/markbates/gentronics"

func New() *gentronics.Generator {
	g := gentronics.New()

	g.Add(gentronics.NewFile(".buffalo.dev.yml", nRefresh))
	return g
}

const nRefresh = `app_root: .
ignored_folders:
- vendor
- log
- logs
- assets
- public
- grifts
- tmp
- bin
- node_modules
- .sass-cache
included_extensions:
- .go
- .html
- .md
- .js
- .tmpl
build_path: tmp
build_delay: 200ns
binary_name: {{name}}-build
command_flags: []
enable_colors: true
log_name: buffalo
`
