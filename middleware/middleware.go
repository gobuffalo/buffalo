package middleware

import "github.com/markbates/oncer"

func init() {
	oncer.Deprecate(0, "github.com/gobuffalo/buffalo/middleware", "Use github.com/gobuffalo/mw-contenttype, github.com/gobuffalo/mw-paramlogger, and github.com/gobuffalo/buffalo-pop instead.")
}
