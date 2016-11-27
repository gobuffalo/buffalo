package buffalo

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/markbates/going/defaults"
)

type Options struct {
	Env            string
	LogLevel       string
	Logger         Logger
	LogDir         string
	NotFound       http.Handler
	MethodOverride func(r *http.Request)
	prefix         string
}

func NewOptions() Options {
	return optionsWithDefaults(Options{})
}

func optionsWithDefaults(opts Options) Options {
	opts.Env = defaults.String(opts.Env, defaults.String(os.Getenv("BUFFALO_ENV"), defaults.String(os.Getenv("GO_ENV"), "development")))
	opts.LogLevel = defaults.String(opts.LogLevel, "debug")
	pwd, _ := os.Getwd()
	opts.LogDir = defaults.String(opts.LogDir, filepath.Join(pwd, "logs"))
	return opts
}
