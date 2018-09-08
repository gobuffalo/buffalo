package middleware

import "github.com/gobuffalo/buffalo-pop/pop/popmw"

// PopTransaction is a piece of Buffalo middleware that wraps each
// request in a transaction. The transaction will automatically get
// committed if there's no errors and the response status code is a
// 2xx or 3xx, otherwise it'll be rolled back. It will also add a
// field to the log, "db", that shows the total duration spent during
// the request making database calls.
//
// Deprecated: use github.com/gobuffalo/buffalo-pop/pop/popmw#Transaction instead.
var PopTransaction = popmw.Transaction
