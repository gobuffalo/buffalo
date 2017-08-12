package plush

import (
	"fmt"
	"html/template"

	"github.com/gobuffalo/tags"
	"github.com/gobuffalo/tags/form"
	"github.com/gobuffalo/tags/form/bootstrap"
)

// FormHelper implements a Plush helper around the
// form.New function in the github.com/gobuffalo/tags/form package
func FormHelper(opts tags.Options, help HelperContext) (template.HTML, error) {
	return helper(opts, help, func(opts tags.Options) helperable {
		return form.New(opts)
	})
}

// FormForHelper implements a Plush helper around the
// form.NewFormFor function in the github.com/gobuffalo/tags/form package
func FormForHelper(model interface{}, opts tags.Options, help HelperContext) (template.HTML, error) {
	return helper(opts, help, func(opts tags.Options) helperable {
		return form.NewFormFor(model, opts)
	})
}

// BootstrapFormHelper implements a Plush helper around the
// bootstrap.New function in the github.com/gobuffalo/tags/form/bootstrap package
func BootstrapFormHelper(opts tags.Options, help HelperContext) (template.HTML, error) {
	return helper(opts, help, func(opts tags.Options) helperable {
		return bootstrap.New(opts)
	})
}

// BootstrapFormForHelper implements a Plush helper around the
// bootstrap.NewFormFor function in the github.com/gobuffalo/tags/form/bootstrap package
func BootstrapFormForHelper(model interface{}, opts tags.Options, help HelperContext) (template.HTML, error) {
	return helper(opts, help, func(opts tags.Options) helperable {
		return bootstrap.NewFormFor(model, opts)
	})
}

type helperable interface {
	SetAuthenticityToken(string)
	Append(...tags.Body)
	HTMLer
}

func helper(opts tags.Options, help HelperContext, fn func(opts tags.Options) helperable) (template.HTML, error) {
	hn := "f"
	if n, ok := opts["var"]; ok {
		hn = n.(string)
		delete(opts, "var")
	}
	if opts["errors"] == nil {
		opts["errors"] = help.Context.Value("errors")
	}
	form := fn(opts)
	if help.Value("authenticity_token") != nil {
		form.SetAuthenticityToken(fmt.Sprint(help.Value("authenticity_token")))
	}
	ctx := help.Context.New()
	ctx.Set(hn, form)
	s, err := help.BlockWith(ctx)
	if err != nil {
		return "", err
	}
	form.Append(s)
	return form.HTML(), nil
}
