package buffalo

import "net/http"

type Options struct {
	Env      string
	LogLevel string
	Logger   Logger
	NotFound http.Handler
	prefix   string
}
