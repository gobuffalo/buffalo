package buffalo

import "net/http"

type Options struct {
	Env            string
	LogLevel       string
	Logger         Logger
	NotFound       http.Handler
	MethodOverride func(r *http.Request)
	prefix         string
}
