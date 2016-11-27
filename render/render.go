package render

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/pkg/errors"
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
	opts.TemplateFuncs["partial"] = func(name string) template.HTML {
		return template.HTML("")
	}
	opts.TemplateFuncs["debug"] = func(data interface{}) template.HTML {
		return template.HTML(fmt.Sprintf("%+v", data))
	}

	opts.templates = template.New("").Funcs(opts.TemplateFuncs)
	if opts.TemplatesPath != "" {
		var err error
		opts.templates, err = parseAndCache(opts.TemplatesPath, opts.TemplateFuncs)
		if err != nil {
			log.Fatal(err)
		}
	}

	e := &Engine{
		Options: opts,
		moot:    &sync.Mutex{},
	}
	return e
}

func parseAndCache(templatesPath string, helpers template.FuncMap) (*template.Template, error) {
	t := template.New("").Funcs(helpers)
	err := filepath.Walk(templatesPath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			b, err := ioutil.ReadFile(path)
			if err != nil {
				return errors.WithStack(err)
			}
			t, err = t.New(info.Name()).Parse(string(b))
			if err != nil {
				return errors.WithStack(err)
			}
		}
		return nil
	})
	return t, errors.WithStack(err)
}
