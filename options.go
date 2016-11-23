package buffalo

import "net/http"

type Options struct {
	Env                string
	LogLevel           string
	Logger             Logger
	DefaultContentType string
	NotFound           http.Handler
	prefix             string
}
