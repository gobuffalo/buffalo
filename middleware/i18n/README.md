i18n middleware
===============

This Buffalo middleware enables i18n features in your app:
* User language detection from configurable sources
* Translation helper using locales bundles from github.com/nicksnyder/go-i18n
* Localized views

Installation
------------

This middleware is setup by default on a new Buffalo app:

**actions/app.go**
```go
var app *buffalo.App

// T is used to provide translations
var T *i18n.Translator

// App is where all routes and middleware for buffalo
// should be defined. This is the nerve center of your
// application.
func App() *buffalo.App {
    if app == nil {
        // [...]

        // Setup and use translations:
	var err error
	if T, err = i18n.New(packr.NewBox("../locales"), "en"); err != nil {
		app.Stop(err)
	}
	app.Use(T.Middleware())
    }
    return app
}
```

Use `i18n.New` to create a new instance of the translation module, then add the middleware (`T.Middleware()`) to the app to enable its features.

See https://gobuffalo.io/docs/localization for further info about Buffalo translation features and configuration.
