module github.com/gobuffalo/buffalo/buffalo

go 1.13

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/cockroachdb/cockroach-go v0.0.0-20190925194419-606b3d062051 // indirect
	github.com/fatih/color v1.7.0
	github.com/gobuffalo/attrs v0.1.0
	github.com/gobuffalo/buffalo v0.15.2-0.20191120211510-584997c3f17d
	github.com/gobuffalo/buffalo-pop v1.23.1
	github.com/gobuffalo/clara v0.9.1
	github.com/gobuffalo/events v1.4.0
	github.com/gobuffalo/fizz v1.9.5 // indirect
	github.com/gobuffalo/flect v0.1.7
	github.com/gobuffalo/genny v0.4.1
	github.com/gobuffalo/licenser v1.4.0
	github.com/gobuffalo/logger v1.0.2
	github.com/gobuffalo/meta v0.2.1
	github.com/gobuffalo/packr/v2 v2.7.1
	github.com/gobuffalo/pop v4.12.2+incompatible
	github.com/jackc/pgconn v1.1.0 // indirect
	github.com/markbates/grift v1.5.0
	github.com/markbates/oncer v1.0.0
	github.com/markbates/refresh v1.8.0
	github.com/markbates/safe v1.0.1
	github.com/markbates/sigtx v1.0.0
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/cobra v0.0.5
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.5.0
	github.com/stretchr/testify v1.4.0
	golang.org/x/sync v0.0.0-20190911185100-cd5d95a43a6e
	golang.org/x/tools v0.0.0-20191124021906-f5828fc9a103
)

replace github.com/gobuffalo/buffalo => ../
