package buffalo

import "encoding/json"

// flashKey is the prefix inside the Session.
const flashKey = "_flash_"

//Flash is a struct that helps with the operations over flash messages.
type Flash struct {
	data map[string][]string
}

//Delete removes a particular key from the Flash.
func (f Flash) Delete(key string) {
	delete(f.data, key)
}

//Clear removes all keys from the Flash.
func (f *Flash) Clear() {
	f.data = map[string][]string{}
}

//Set allows to set a list of values into a particular key.
func (f Flash) Set(key string, values []string) {
	f.data[key] = values
}

//Add adds a flash value for a flash key, if the key already has values the list for that value grows.
func (f Flash) Add(key, value string) {
	if len(f.data[key]) == 0 {
		f.data[key] = []string{value}
		return
	}

	f.data[key] = append(f.data[key], value)
}

//Persist the flash inside the session.
func (f Flash) persist(session *Session) {
	b, _ := json.Marshal(f.data)
	session.Set(flashKey, b)
}

//newFlash creates a new Flash and loads the session data inside its data.
func newFlash(session *Session) *Flash {
	result := &Flash{
		data: map[string][]string{},
	}

	if session.Session != nil {
		if f := session.Get(flashKey); f != nil {
			json.Unmarshal(f.([]byte), &result.data)
		}
	}
	return result
}
