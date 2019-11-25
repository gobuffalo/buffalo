package consts

const (
	// Environments
	Development = "development"
	Production  = "production"
	Test        = "test"

	// Defaults
	Def_Addr          = "0.0.0.0"
	Def_AddrDev       = "127.0.0.1"
	Def_AssestsMaxAge = "31536000"
	Def_Port          = "3000"
	Def_Root          = "/"
	Def_SessionName   = "_buffalo_session"

	// HTTP Verbs
	HTTP_DELETE   = "DELETE"
	HTTP_GET      = "GET"
	HTTP_HEAD     = "HEAD"
	HTTP_OPTIONS  = "OPTIONS"
	HTTP_PATCH    = "PATCH"
	HTTP_POST     = "POST"
	HTTP_PUT      = "PUT"
	HTTP_Override = "_method"

	// HTTP
	HTTP_CORS               = "Access-Control-Allow-Origin"
	HTTP_CacheControl       = "Cache-Control"
	HTTP_Connection         = "Connection"
	HTTP_ContentDisposition = "Content-Disposition"
	HTTP_ContentLength      = "Content-Length"
	HTTP_ContentType        = "Content-Type"
	HTTP_ETag               = "ETag"

	// MIME Types
	MIME_HTML         = "text/html; charset=utf-8"
	MIME_JSON         = "application/json; charset=utf-8"
	MIME_JavaScript   = "application/javascript; charset=utf-8"
	MIME_Octet_Stream = "application/octet-stream"
	MIME_Text         = "text/plain; charset=utf-8"
	MIME_XML          = "application/xml; charset=utf-8"

	// Requests
	REQ_ID          = "request_id"
	REQ_RequestorID = "requestor_id"

	// Environment Variables
	ADDR           = "ADDR"
	ASSETS_MAX_AGE = "ASSETS_MAX_AGE"
	GO_ENV         = "GO_ENV"
	HOST           = "HOST"
	LOG_LEVEL      = "LOG_LEVEL"
	PORT           = "PORT"
	SESSION_SECRET = "SESSION_SECRET"

	BUFFALO_PLUGIN_CACHE = "BUFFALO_PLUGIN_CACHE"

	// Other

	NL = "\n"
)
