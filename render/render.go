package render

import (
	"fmt"
	"html/template"
	"sync"
)

type Engine struct {
	*Options
	moot *sync.Mutex
}

func New(opts *Options) *Engine {
	if opts.TemplateFuncs == nil {
		opts.TemplateFuncs = template.FuncMap{}
	}
	opts.TemplateFuncs["yield"] = func() template.HTML {
		return template.HTML("")
	}
	opts.TemplateFuncs["debug"] = func(data interface{}) template.HTML {
		return template.HTML(fmt.Sprintf("%+v", data))
	}
	opts.templates = template.New("").Funcs(opts.TemplateFuncs)
	e := &Engine{
		Options: opts,
		moot:    &sync.Mutex{},
	}
	return e
}
