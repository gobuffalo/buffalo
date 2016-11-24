package buffalo

import (
	"net/http"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/markbates/going/defaults"
)

type Options struct {
	Env            string
	LogLevel       string
	Logger         Logger
	NotFound       http.Handler
	MethodOverride func(r *http.Request)
	prefix         string
}

func NewOptions() Options {
	return optionsWithDefaults(Options{})
}

func optionsWithDefaults(opts Options) Options {
	opts.Env = defaults.String(opts.Env, defaults.String(os.Getenv("BUFFALO_ENV"), defaults.String(os.Getenv("GO_ENV"), "development")))
	if opts.Logger == nil {
		l := logrus.New()
		l.Level, _ = logrus.ParseLevel(defaults.String(opts.LogLevel, "debug"))
		opts.Logger = l
	}
	return opts
}
