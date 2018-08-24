package middleware

import (
	contenttype "github.com/gobuffalo/mw-contenttype"
)

// SetContentType on the request to desired type. This will
// override any content type sent by the client.
//
// Deprecated: use github.com/gobuffalo/mw-contenttype#Set instead.
var SetContentType = contenttype.Set

// AddContentType will add a secondary content type to
// a request. If no content type is sent by the client
// the default will be set, otherwise the client's
// content type will be used.
//
// Deprecated: use github.com/gobuffalo/mw-contenttype#Add instead.
var AddContentType = contenttype.Add
