package buffalo

import (
	"fmt"
	"strings"
)

//Flash is a struct that helps with the operations over flash messages.
type Flash struct {
	Session *Session
}

//Set sets a message inside the Flash.
func (f *Flash) Set(key, value string) {
	f.Session.Set(fmt.Sprintf("_flash_%v", key), value)
	f.Session.Save()
}

//Get gets a message from inside the Flash.
func (f *Flash) Get(key string) string {
	val := f.Session.Get(fmt.Sprintf("_flash_%v", key))
	if val == nil {
		return ""
	}

	return val.(string)
}

//Delete removes a particular key from the Flash.
func (f *Flash) Delete(key string) {
	f.Session.Delete(fmt.Sprintf("_flash_%v", key))
	f.Session.Save()
}

//Clear Wipes all the flash messages.
func (f *Flash) Clear() {
	for k := range f.Session.Session.Values {
		if strings.HasPrefix(k.(string), "_flash_") {
			f.Session.Delete(k)
		}
	}
	f.Session.Save()
}

//Creates a new flash instance with the passed session and empty data
func newFlash(session *Session) *Flash {
	flash := &Flash{
		Session: session,
	}

	return flash
}
