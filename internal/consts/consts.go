package consts

const (
	// Environments
	Development = "development"
	Production  = "production"
	Test        = "test"

	// Defaults
	Def_Port          = "3000"
	Def_Addr          = "0.0.0.0"
	Def_AddrDev       = "127.0.0.1"
	Def_SessionName   = "_buffalo_session"
	Def_Root          = "/"
	Def_AssestsMaxAge = "31536000"

	// HTTP Verbs
	HTTP_DELETE  = "DELETE"
	HTTP_GET     = "GET"
	HTTP_HEAD    = "HEAD"
	HTTP_OPTIONS = "OPTIONS"
	HTTP_PATCH   = "PATCH"
	HTTP_POST    = "POST"
	HTTP_PUT     = "PUT"

	// HTTP
	HTTP_ETag         = "ETag"
	HTTP_CacheControl = "Cache-Control"

	// Environment Variables
	ADDR           = "ADDR"
	GO_ENV         = "GO_ENV"
	HOST           = "HOST"
	LOG_LEVEL      = "LOG_LEVEL"
	PORT           = "PORT"
	SESSION_SECRET = "SESSION_SECRET"
	ASSETS_MAX_AGE = "ASSETS_MAX_AGE"
)
