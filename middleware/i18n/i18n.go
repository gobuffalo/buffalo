package i18n

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/packr"
	"github.com/nicksnyder/go-i18n/i18n"
	"github.com/nicksnyder/go-i18n/i18n/language"
	"github.com/nicksnyder/go-i18n/i18n/translation"
	"github.com/pkg/errors"
)

// LanguageExtractor can be implemented for custom finding of search
// languages. This can be useful if you want to load a user's language
// from something like a database. See Middleware() for more information
// on how the default implementation searches for languages.
//
// Deprecated: use github.com/gobuffalo/mw-i18n#LanguageExtractor instead.
type LanguageExtractor func(LanguageExtractorOptions, buffalo.Context) []string

// LanguageExtractorOptions is a map of options for a LanguageExtractor.
//
// Deprecated: use github.com/gobuffalo/mw-i18n#LanguageExtractorOptions instead.
type LanguageExtractorOptions map[string]interface{}

// Translator for handling all your i18n needs.
//
// Deprecated: use github.com/gobuffalo/mw-i18n#Translator instead.
type Translator struct {
	// Box - where are the files?
	Box packr.Box
	// DefaultLanguage - default is passed as a parameter on New.
	DefaultLanguage string
	// HelperName - name of the view helper. default is "t"
	HelperName string
	// LanguageExtractors - a sorted list of user language extractors.
	LanguageExtractors []LanguageExtractor
	// LanguageExtractorOptions - a map with options to give to LanguageExtractors.
	LanguageExtractorOptions LanguageExtractorOptions
}

// Load translations from the t.Box.
//
// Deprecated: use github.com/gobuffalo/mw-i18n#Load instead.
func (t *Translator) Load() error {
	fmt.Printf("i18n.Load is deprecated and will be removed in the next version. Please use github.com/gobuffalo/mw-i18n#Load instead.")
	return t.Box.Walk(func(path string, f packr.File) error {
		b, err := t.Box.MustBytes(path)
		if err != nil {
			return errors.Wrapf(err, "unable to read locale file %s", path)
		}

		base := filepath.Base(path)
		dir := filepath.Dir(path)

		// Add a prefix to the loaded string, to avoid collision with an ISO lang code
		err = i18n.ParseTranslationFileBytes(fmt.Sprintf("%sbuff%s", dir, base), b)
		if err != nil {
			return errors.Wrapf(err, "unable to parse locale file %s", base)
		}
		return nil
	})
}

// AddTranslation directly, without using a file. This is useful if you wish to load translations
// from a database, instead of disk.
//
// Deprecated: use github.com/gobuffalo/mw-i18n#AddTranslation instead.
func (t *Translator) AddTranslation(lang *language.Language, translations ...translation.Translation) {
	fmt.Printf("i18n.AddTranslation is deprecated and will be removed in the next version. Please use github.com/gobuffalo/mw-i18n#AddTranslation instead.")
	i18n.AddTranslation(lang, translations...)
}

// New Translator. Requires a packr.Box that points to the location
// of the translation files, as well as a default language. This will
// also call t.Load() and load the translations from disk.
//
// Deprecated: use github.com/gobuffalo/mw-i18n#New instead.
func New(box packr.Box, language string) (*Translator, error) {
	t := &Translator{
		Box:             box,
		DefaultLanguage: language,
		HelperName:      "t",
		LanguageExtractorOptions: LanguageExtractorOptions{
			"CookieName":    "lang",
			"SessionName":   "lang",
			"URLPrefixName": "lang",
		},
		LanguageExtractors: []LanguageExtractor{
			CookieLanguageExtractor,
			SessionLanguageExtractor,
			HeaderLanguageExtractor,
		},
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
//
// Deprecated: use github.com/gobuffalo/mw-i18n#Middleware instead.
func (t *Translator) Middleware() buffalo.MiddlewareFunc {
	return func(next buffalo.Handler) buffalo.Handler {
		return func(c buffalo.Context) error {
			fmt.Printf("i18n.Middleware is deprecated and will be removed in the next version. Please use github.com/gobuffalo/mw-i18n#Middleware instead.")

			// in development reload the translations
			if c.Value("env").(string) == "development" {
				err := t.Load()
				if err != nil {
					return err
				}
			}

			// set languages in context, if not set yet
			if langs := c.Value("languages"); langs == nil {
				c.Set("languages", t.extractLanguage(c))
			}

			// set translator
			if T := c.Value("T"); T == nil {
				langs := c.Value("languages").([]string)
				T, err := i18n.Tfunc(langs[0], langs[1:]...)
				if err != nil {
					c.Logger().Warn(err)
					c.Logger().Warn("Your locale files are probably empty or missing")
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
//
// Deprecated: use github.com/gobuffalo/mw-i18n#Translate instead.
func (t *Translator) Translate(c buffalo.Context, translationID string, args ...interface{}) string {
	fmt.Printf("i18n.Translate is deprecated and will be removed in the next version. Please use github.com/gobuffalo/mw-i18n#Translate instead.")
	T := c.Value("T").(i18n.TranslateFunc)
	return T(translationID, args...)
}

// AvailableLanguages gets the list of languages provided by the app.
//
// Deprecated: use github.com/gobuffalo/mw-i18n#AvailableLanguages instead.
func (t *Translator) AvailableLanguages() []string {
	fmt.Printf("i18n.AvailableLanguages is deprecated and will be removed in the next version. Please use github.com/gobuffalo/mw-i18n#AvailableLanguages instead.")
	lt := i18n.LanguageTags()
	sort.Strings(lt)
	return lt
}

// Refresh updates the context, reloading translation functions.
// It can be used after language change, to be able to use translation functions
// in the new language (for a flash message, for instance).
//
// Deprecated: use github.com/gobuffalo/mw-i18n#Refresh instead.
func (t *Translator) Refresh(c buffalo.Context, newLang string) {
	fmt.Printf("i18n.Refresh is deprecated and will be removed in the next version. Please use github.com/gobuffalo/mw-i18n instead.")
	langs := []string{newLang}
	langs = append(langs, t.extractLanguage(c)...)

	// Refresh languages
	c.Set("languages", langs)

	T, err := i18n.Tfunc(langs[0], langs[1:]...)
	if err != nil {
		c.Logger().Warn(err)
		c.Logger().Warn("Your locale files are probably empty or missing")
	}

	// Refresh translation engine
	c.Set("T", T)
}

func (t *Translator) extractLanguage(c buffalo.Context) []string {
	langs := []string{}
	for _, extractor := range t.LanguageExtractors {
		langs = append(langs, extractor(t.LanguageExtractorOptions, c)...)
	}
	// Add default language, even if no language extractor is defined
	langs = append(langs, t.DefaultLanguage)
	return langs
}

// CookieLanguageExtractor is a LanguageExtractor implementation, using a cookie.
//
// Deprecated: use github.com/gobuffalo/mw-i18n#CookieLanguageExtractor instead.
func CookieLanguageExtractor(o LanguageExtractorOptions, c buffalo.Context) []string {
	fmt.Printf("CookieLanguageExtractor is deprecated and will be removed in the next version. Please use github.com/gobuffalo/mw-i18n#CookieLanguageExtractor instead.")
	langs := make([]string, 0)
	// try to get the language from a cookie:
	if cookieName := o["CookieName"].(string); cookieName != "" {
		if cookie, err := c.Request().Cookie(cookieName); err == nil {
			if cookie.Value != "" {
				langs = append(langs, cookie.Value)
			}
		}
	} else {
		c.Logger().Error("i18n middleware: \"CookieName\" is not defined in LanguageExtractorOptions")
	}
	return langs
}

// SessionLanguageExtractor is a LanguageExtractor implementation, using a session.
//
// Deprecated: use github.com/gobuffalo/mw-i18n#SessionLanguageExtractor instead.
func SessionLanguageExtractor(o LanguageExtractorOptions, c buffalo.Context) []string {
	fmt.Printf("SessionLanguageExtractor is deprecated and will be removed in the next version. Please use github.com/gobuffalo/mw-i18n#SessionLanguageExtractor instead.")
	langs := make([]string, 0)
	// try to get the language from the session
	if sessionName := o["SessionName"].(string); sessionName != "" {
		if s := c.Session().Get(sessionName); s != nil {
			langs = append(langs, s.(string))
		}
	} else {
		c.Logger().Error("i18n middleware: \"SessionName\" is not defined in LanguageExtractorOptions")
	}
	return langs
}

// HeaderLanguageExtractor is a LanguageExtractor implementation, using a HTTP Accept-Language
// header.
//
// Deprecated: use github.com/gobuffalo/mw-i18n#HeaderLanguageExtractor instead.
func HeaderLanguageExtractor(o LanguageExtractorOptions, c buffalo.Context) []string {
	fmt.Printf("HeaderLanguageExtractor is deprecated and will be removed in the next version. Please use github.com/gobuffalo/mw-i18n#HeaderLanguageExtractor instead.")
	langs := make([]string, 0)
	// try to get the language from a header:
	acceptLang := c.Request().Header.Get("Accept-Language")
	if acceptLang != "" {
		langs = append(langs, parseAcceptLanguage(acceptLang)...)
	}
	return langs
}

// URLPrefixLanguageExtractor is a LanguageExtractor implementation, using a prefix in the URL.
//
// Deprecated: use github.com/gobuffalo/mw-i18n#URLPrefixLanguageExtractor instead.
func URLPrefixLanguageExtractor(o LanguageExtractorOptions, c buffalo.Context) []string {
	fmt.Printf("URLPrefixLanguageExtractor is deprecated and will be removed in the next version. Please use github.com/gobuffalo/mw-i18n#URLPrefixLanguageExtractor instead.")
	langs := make([]string, 0)
	// try to get the language from an URL prefix:
	if urlPrefixName := o["URLPrefixName"].(string); urlPrefixName != "" {
		paramLang := c.Param(urlPrefixName)
		if paramLang != "" && strings.HasPrefix(c.Request().URL.Path, fmt.Sprintf("/%s", paramLang)) {
			langs = append(langs, paramLang)
		}
	} else {
		c.Logger().Error("i18n middleware: \"URLPrefixName\" is not defined in LanguageExtractorOptions")
	}
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
