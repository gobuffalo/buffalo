package i18n

import (
	"log"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/packr"
	"github.com/nicksnyder/go-i18n/i18n"
	"github.com/nicksnyder/go-i18n/i18n/language"
	"github.com/nicksnyder/go-i18n/i18n/translation"
	"github.com/pkg/errors"
)

// LanguageFinder can be implemented for custom finding of search
// languages. This can be useful if you want to load a user's langugage
// from something like a database. See Middleware() for more information
// on how the default implementation searches for languages.
type LanguageFinder func(*Translator, buffalo.Context) []string

// Translator for handling all your i18n needs.
type Translator struct {
	// Box - where are the files?
	Box packr.Box
	// DefaultLanguage - default is "en-US"
	DefaultLanguage string
	// CookieName - name of the cookie to find the desired language.
	// default is "lang"
	CookieName string
	// SessionName - name of the session to find the desired language.
	// default is "lang"
	SessionName string
	// HelperName - name of the view helper. default is "t"
	HelperName     string
	LanguageFinder LanguageFinder
}

// Load translations from the t.Box.
func (t *Translator) Load() error {
	return t.Box.Walk(func(path string, f packr.File) error {
		b, err := t.Box.MustBytes(path)
		if err != nil {
			log.Fatal(err)
			return errors.WithStack(err)
		}
		return i18n.ParseTranslationFileBytes(path, b)
	})
}

// AddTranslation directly, without using a file. This is useful if you wish to load translations
// from a database, instead of disk.
func (t *Translator) AddTranslation(lang *language.Language, translations ...translation.Translation) {
	i18n.AddTranslation(lang, translations...)
}

// New Translator. Requires a packr.Box that points to the location
// of the translation files, as well as a default language. This will
// also call t.Load() and load the translations from disk.
func New(box packr.Box, language string) (*Translator, error) {
	t := &Translator{
		Box:             box,
		DefaultLanguage: language,
		CookieName:      "lang",
		SessionName:     "lang",
		HelperName:      "t",
		LanguageFinder:  defaultLanguageFinder,
	}
	return t, t.Load()
}

// Middleware for loading the translations for the language(s)
// selected. By default languages are loaded in the following order:
//
// Cookie - "lang"
// Header - "Accept-Language"
// Default - "en-US"
//
// These values can be changed on the Translator itself. In development
// model the translation files will be reloaded on each request.
func (t *Translator) Middleware() buffalo.MiddlewareFunc {
	return func(next buffalo.Handler) buffalo.Handler {
		return func(c buffalo.Context) error {

			// in development reload the translations
			if c.Value("env").(string) == "development" {
				err := t.Load()
				if err != nil {
					return err
				}
			}

			// set up the helper function for the views:
			c.Set(t.HelperName, func(s string) (string, error) {
				return t.Translate(c, s)
			})
			return next(c)
		}
	}
}

// Translate a string given a Context
func (t *Translator) Translate(c buffalo.Context, s string) (string, error) {
	if langs := c.Value("languages"); langs == nil {
		c.Set("languages", t.LanguageFinder(t, c))
	}
	langs := c.Value("languages").([]string)
	T, err := i18n.Tfunc(langs[0], langs[1:]...)
	if err != nil {
		return "", err
	}
	return T(s, c.Data()), nil
}

func defaultLanguageFinder(t *Translator, c buffalo.Context) []string {
	langs := []string{}

	r := c.Request()

	// try to get the language from a cookie:
	if cookie, err := r.Cookie(t.CookieName); err == nil {
		if cookie.Value != "" {
			langs = append(langs, cookie.Value)
		}
	}

	// try to get the language from the session
	if s := c.Session().Get(t.SessionName); s != nil {
		langs = append(langs, s.(string))
	}

	// try to get the language from a header:
	acceptLang := r.Header.Get("Accept-Language")
	if acceptLang != "" {
		langs = append(langs, acceptLang)
	}

	// try to get the language from the session:
	langs = append(langs, t.DefaultLanguage)
	return langs
}
