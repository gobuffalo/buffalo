package buffalo

import (
	"net/http"
	"os"

	"github.com/gorilla/sessions"
)

// Store is the `github.com/gorilla/sessions` store used to back
// the session. It defaults to use a cookie store and the ENV variable
// `SESSION_SECRET`.
var SessionStore = sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))

// SessionName is the name of the session cookie that is set. This defaults
// to "_buffalo_session".
var SessionName = "_buffalo_session"

type Session struct {
	Session *sessions.Session
	req     *http.Request
	res     http.ResponseWriter
}

// Save the current session
func (s *Session) Save() error {
	return s.Session.Save(s.req, s.res)
}

// Get a value from the current session
func (s *Session) Get(name interface{}) interface{} {
	return s.Session.Values[name]
}

// Set a value onto the current session. If a value with that name
// already exists it will be overridden with the new value.
func (s *Session) Set(name, value interface{}) {
	s.Session.Values[name] = value
}

// Delete a value from the current session.
func (s *Session) Delete(name interface{}) {
	delete(s.Session.Values, name)
}

// Get a session using a request and response.
func GetSession(r *http.Request, w http.ResponseWriter) *Session {
	session, _ := SessionStore.Get(r, SessionName)
	return &Session{
		Session: session,
		req:     r,
		res:     w,
	}
}
