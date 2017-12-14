package i18n

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/packr"
	"github.com/nicksnyder/go-i18n/i18n"
	"github.com/nicksnyder/go-i18n/i18n/language"
	"github.com/nicksnyder/go-i18n/i18n/translation"
	"github.com/pkg/errors"
)

// LanguageFinder can be implemented for custom finding of search
// languages. This can be useful if you want to load a user's language
// from something like a database. See Middleware() for more information
// on how the default implementation searches for languages.
type LanguageFinder func(*Translator, buffalo.Context) []string

// Translator for handling all your i18n needs.
type Translator struct {
	// Box - where are the files?
	Box packr.Box
	// DefaultLanguage - default is passed as a parameter on New.
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

		base := filepath.Base(path)
		dir := filepath.Dir(path)

		// Add a prefix to the loaded string, to avoid collision with an ISO lang code
		return i18n.ParseTranslationFileBytes(fmt.Sprintf("%sbuff%s", dir, base), b)
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
// Session - "lang"
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

			// set languages in context, if not set yet
			if langs := c.Value("languages"); langs == nil {
				c.Set("languages", t.LanguageFinder(t, c))
			}

			// set translator
			if T := c.Value("T"); T == nil {
				langs := c.Value("languages").([]string)
				T, err := i18n.Tfunc(langs[0], langs[1:]...)
				if err != nil {
					return err
				}
				c.Set("T", T)
			}

			// set up the helper function for the views:
			c.Set(t.HelperName, func(s string, i ...interface{}) string {
				return t.Translate(c, s, i...)
			})
			return next(c)
		}
	}
}

// Translate returns the translation of the string identified by translationID.
//
// See https://github.com/nicksnyder/go-i18n
//
// If there is no translation for translationID, then the translationID itself is returned.
// This makes it easy to identify missing translations in your app.
//
// If translationID is a non-plural form, then the first variadic argument may be a map[string]interface{}
// or struct that contains template data.
//
// If translationID is a plural form, the function accepts two parameter signatures
// 1. T(count int, data struct{})
// The first variadic argument must be an integer type
// (int, int8, int16, int32, int64) or a float formatted as a string (e.g. "123.45").
// The second variadic argument may be a map[string]interface{} or struct{} that contains template data.
// 2. T(data struct{})
// data must be a struct{} or map[string]interface{} that contains a Count field and the template data,
// Count field must be an integer type (int, int8, int16, int32, int64)
// or a float formatted as a string (e.g. "123.45").
func (t *Translator) Translate(c buffalo.Context, translationID string, args ...interface{}) string {
	T := c.Value("T").(i18n.TranslateFunc)
	return T(translationID, args...)
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
		langs = append(langs, parseAcceptLanguage(acceptLang)...)
	}

	// finally set the default app language as fallback
	langs = append(langs, t.DefaultLanguage)
	return langs
}

// Inspired from https://siongui.github.io/2015/02/22/go-parse-accept-language/
// Parse an Accept-Language string to get usable lang values for i18n system
func parseAcceptLanguage(acptLang string) []string {
	var lqs []string

	langQStrs := strings.Split(acptLang, ",")
	for _, langQStr := range langQStrs {
		trimedLangQStr := strings.Trim(langQStr, " ")

		langQ := strings.Split(trimedLangQStr, ";")
		lq := langQ[0]
		lqs = append(lqs, lq)
	}
	return lqs
}
