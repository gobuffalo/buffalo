package buffalo

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

// FlashPrefix is the prefix inside the Session.
const FlashPrefix = "_flash_"

//Flash is a struct that helps with the operations over flash messages.
type Flash struct {
	data map[string][]string
}

//Set sets a message inside the Flash.
func (f *Flash) Set(key, value string) {
	f.data[key] = []string{value}
}

//Get gets a message from inside the Flash.
func (f *Flash) Get(key string) []string {
	return f.data[key]
}

//Delete removes a particular key from the Flash.
func (f *Flash) Delete(key string) {
	delete(f.data, key)
}

//Add adds a flash value for a flash key, if the key already has values the list for that value grows.
func (f *Flash) Add(key, value string) {
	if len(f.data[key]) == 0 {
		f.data[key] = []string{value}
		return
	}

	f.data[key] = append(f.data[key], value)
}

//Clear Wipes all the flash messages.
func (f *Flash) Clear() {
	f.data = map[string][]string{}
}

//Persist the flash inside the session.
func (f *Flash) Persist(session *Session) {
	for k, v := range f.data {
		sessionKey := fmt.Sprintf("%v%v", FlashPrefix, k)
		bson, err := json.Marshal(v)

		if err == nil {
			session.Set(sessionKey, string(bson))
		}
	}

	session.Save()
}

//newFlash creates a new Flash and loads the session data inside its data.
func newFlash(session *Session) *Flash {
	result := &Flash{
		data: map[string][]string{},
	}

	if session.Session != nil {
		for k := range session.Session.Values {
			sessionName := k.(string)
			if strings.HasPrefix(sessionName, FlashPrefix) {
				log.Println(k.(string))
				flashName := strings.Replace(sessionName, FlashPrefix, "", -1)

				var flashes []string
				err := json.Unmarshal([]byte(session.Get(sessionName).(string)), &flashes)
				if err == nil {
					result.data[flashName] = flashes
				}
			}
		}
	}

	return result
}
