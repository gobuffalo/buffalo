package middleware

import (
	"github.com/gobuffalo/buffalo"
	paramlogger "github.com/gobuffalo/mw-paramlogger"
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
	// Ensure the exclusion list is forwarded
	paramlogger.ParameterExclusionList = ParameterExclusionList
	return paramlogger.ParameterLogger(next)
}
