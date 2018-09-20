package i18n

import (
	i18nm "github.com/gobuffalo/mw-i18n"
	"github.com/markbates/oncer"
)

// LanguageExtractor can be implemented for custom finding of search
// languages. This can be useful if you want to load a user's language
// from something like a database. See Middleware() for more information
// on how the default implementation searches for languages.
//
// Deprecated: use github.com/gobuffalo/mw-i18n#LanguageExtractor instead.
type LanguageExtractor = i18nm.LanguageExtractor

// LanguageExtractorOptions is a map of options for a LanguageExtractor.
//
// Deprecated: use github.com/gobuffalo/mw-i18n#LanguageExtractorOptions instead.
type LanguageExtractorOptions = i18nm.LanguageExtractorOptions

// Translator for handling all your i18n needs.
//
// Deprecated: use github.com/gobuffalo/mw-i18n#Translator instead.
type Translator = i18nm.Translator

// CookieLanguageExtractor is a LanguageExtractor implementation, using a cookie.
//
// Deprecated: use github.com/gobuffalo/mw-i18n#CookieLanguageExtractor instead.
var CookieLanguageExtractor = i18nm.CookieLanguageExtractor

// SessionLanguageExtractor is a LanguageExtractor implementation, using a session.
//
// Deprecated: use github.com/gobuffalo/mw-i18n#SessionLanguageExtractor instead.
var SessionLanguageExtractor = i18nm.SessionLanguageExtractor

// HeaderLanguageExtractor is a LanguageExtractor implementation, using a HTTP Accept-Language
// header.
//
// Deprecated: use github.com/gobuffalo/mw-i18n#HeaderLanguageExtractor instead.
var HeaderLanguageExtractor = i18nm.HeaderLanguageExtractor

// URLPrefixLanguageExtractor is a LanguageExtractor implementation, using a prefix in the URL.
//
// Deprecated: use github.com/gobuffalo/mw-i18n#URLPrefixLanguageExtractor instead.
var URLPrefixLanguageExtractor = i18nm.URLPrefixLanguageExtractor

// New Translator. Requires a packr.Box that points to the location
// of the translation files, as well as a default language. This will
// also call t.Load() and load the translations from disk.
//
// Deprecated: use github.com/gobuffalo/mw-i18n#New instead.
var New = i18nm.New

func init() {
	oncer.Deprecate(0, "github.com/gobuffalo/buffalo/middleware/i18n", "Use github.com/gobuffalo/mw-i18n instead.")
}
