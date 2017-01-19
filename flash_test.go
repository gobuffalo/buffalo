package buffalo

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_FlashSet(t *testing.T) {
	r := require.New(t)
	req, _ := http.NewRequest("GET", "/", nil)
	res := httptest.NewRecorder()
	a := Automatic(Options{})

	session := a.getSession(req, res)
	c := DefaultContext{
		session: session,
		flash:   newFlash(session),
	}

	c.Flash().Set("error", "Error!")
	r.Equal(c.Flash().Get("error"), "Error!")
}

func Test_Flash(t *testing.T) {
	r := require.New(t)
	req, _ := http.NewRequest("GET", "/", nil)
	res := httptest.NewRecorder()
	a := Automatic(Options{})

	session := a.getSession(req, res)
	c := DefaultContext{
		session: session,
		flash:   newFlash(session),
	}

	c.Flash().Set("error", "error")
	c.Flash().Set("message", "message")
	c.Flash().Set("warning", "warning")

	session = a.getSession(req, res)
	r.Equal(session.Get("_flash_error").(string), "error")
	r.Equal(c.Flash().Get("error"), "error")

	c.Flash().Delete("error")

	session = a.getSession(req, res)

	r.Equal(session.Get("_flash_error"), nil)
	r.Equal(session.Get("_flash_message"), "message")
	r.Equal(session.Get("_flash_warning"), "warning")
}

func Test_FlashClear(t *testing.T) {
	r := require.New(t)
	req, _ := http.NewRequest("GET", "/", nil)
	res := httptest.NewRecorder()
	a := Automatic(Options{})

	session := a.getSession(req, res)
	c := DefaultContext{
		session: session,
		flash:   newFlash(session),
	}

	c.Flash().Set("error", "error")
	c.Flash().Set("message", "message")
	c.Flash().Set("warning", "warning")

	session = a.getSession(req, res)

	session = a.getSession(req, res)
	r.Equal(session.Get("_flash_error").(string), "error")
	r.Equal(session.Get("_flash_message").(string), "message")
	r.Equal(session.Get("_flash_warning").(string), "warning")

	r.Equal(c.Flash().Get("error"), "error")
	r.Equal(c.Flash().Get("warning"), "warning")
	r.Equal(c.Flash().Get("message"), "message")

	c.Flash().Clear()

	session = a.getSession(req, res)

	r.Equal(session.Get("_flash_error"), nil)
	r.Equal(session.Get("_flash_message"), nil)
	r.Equal(session.Get("_flash_warning"), nil)

	r.Equal(c.Flash().Get("error"), "")
	r.Equal(c.Flash().Get("message"), "")
	r.Equal(c.Flash().Get("warning"), "")

}
