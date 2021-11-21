module github.com/gobuffalo/buffalo

go 1.16

require (
	github.com/BurntSushi/toml v0.4.1
	github.com/dustin/go-humanize v1.0.0
	github.com/fatih/color v1.13.0
	github.com/gobuffalo/envy v1.10.1
	github.com/gobuffalo/events v1.4.2
	github.com/gobuffalo/flect v0.2.4
	github.com/gobuffalo/github_flavored_markdown v1.1.1
	github.com/gobuffalo/helpers v0.6.4
	github.com/gobuffalo/httptest v1.5.1
	github.com/gobuffalo/logger v1.0.6
	github.com/gobuffalo/meta v0.3.1
	github.com/gobuffalo/nulls v0.4.1
	github.com/gobuffalo/plush/v4 v4.1.9
	github.com/gobuffalo/pop/v6 v6.0.0
	github.com/gobuffalo/tags/v3 v3.1.2
	github.com/gorilla/handlers v1.5.1
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/sessions v1.2.1
	github.com/karrick/godirwalk v1.16.1
	github.com/markbates/grift v1.5.0
	github.com/markbates/oncer v1.0.0
	github.com/markbates/refresh v1.12.0
	github.com/markbates/safe v1.0.1
	github.com/markbates/sigtx v1.0.0
	github.com/monoculum/formam v3.5.5+incompatible
	github.com/psanford/memfs v0.0.0-20210214183328-a001468d78ef
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/cobra v1.2.1
	github.com/stretchr/testify v1.7.0
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc
)

replace github.com/gobuffalo/pop/v6 v6.0.0 => github.com/fasmat/pop/v6 v6.0.0-20211121195140-6d95c111f911
