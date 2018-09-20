package middleware

import (
	"github.com/gobuffalo/buffalo"
	paramlogger "github.com/gobuffalo/mw-paramlogger"
	"github.com/markbates/oncer"
)

// ParameterExclusionList is the list of parameter names that will be filtered
// from the application logs (see maskSecrets).
// Important: this list will be used in case insensitive.
//
// Deprecated: use github.com/gobuffalo/mw-paramlogger#ParameterExclusionList instead.
var ParameterExclusionList = paramlogger.ParameterExclusionList

// ParameterLogger logs form and parameter values to the loggers
//
// Deprecated: use github.com/gobuffalo/mw-paramlogger#ParameterLogger instead.
func ParameterLogger(next buffalo.Handler) buffalo.Handler {
	oncer.Deprecate(0, "middleware.ParameterLogger", "Use github.com/gobuffalo/mw-paramlogger.ParameterLogger instead.")
	// Ensure the exclusion list is forwarded
	paramlogger.ParameterExclusionList = ParameterExclusionList
	return paramlogger.ParameterLogger(next)
}
