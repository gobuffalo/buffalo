package buffalo

import (
	"net/http"
	"time"
)

// Cookies allows you to easily get cookies from the request, and set cookies on the response.
type Cookies struct {
	req *http.Request
	res http.ResponseWriter
}

// Get returns the value of the cookie with the given name. Returns http.ErrNoCookie if there's no cookie with that name in the request.
func (c *Cookies) Get(name string) (string, error) {
	ck, err := c.req.Cookie(name)
	if err != nil {
		return "", err
	}

	return ck.Value, nil
}

// Set a cookie on the response, which will expire after the given duration.
func (c *Cookies) Set(name, value string, maxAge time.Duration) {
	ck := http.Cookie{
		Name:   name,
		Value:  value,
		MaxAge: int(maxAge.Seconds()),
	}

	http.SetCookie(c.res, &ck)
}

// SetWithExpirationTime sets a cookie that will expire at a specific time.
// Note that the time is determined by the client's browser, so it might not expire at the expected time,
// for example if the client has changed the time on their computer.
func (c *Cookies) SetWithExpirationTime(name, value string, expires time.Time) {
	ck := http.Cookie{
		Name:    name,
		Value:   value,
		Expires: expires,
	}

	http.SetCookie(c.res, &ck)
}

// Delete sets a header that tells the browser to remove the cookie with the given name.
func (c *Cookies) Delete(name string) {
	ck := http.Cookie{
		Name:  name,
		Value: "v",
		// Setting a time in the distant past, like the unix epoch, removes the cookie,
		// since it has long expired.
		Expires: time.Unix(0, 0),
	}

	http.SetCookie(c.res, &ck)
}
