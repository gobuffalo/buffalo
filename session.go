package buffalo

import (
	"net/http"

	"github.com/gorilla/sessions"
)

// Session wraps the "github.com/gorilla/sessions" API
// in something a little cleaner and a bit more useable.
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
func (a *App) getSession(r *http.Request, w http.ResponseWriter) *Session {
	session, _ := a.SessionStore.Get(r, a.SessionName)
	return &Session{
		Session: session,
		req:     r,
		res:     w,
	}
}
