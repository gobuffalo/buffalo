package core

import "github.com/gobuffalo/buffalo-cli/genny/newapp/core"

// ErrNotInGoPath is thrown when not using go modules outside of GOPATH
var ErrNotInGoPath = core.ErrNotInGoPath

var New = core.New

type Options = core.Options
