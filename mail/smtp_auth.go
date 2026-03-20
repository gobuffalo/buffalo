// Portions of this code are derived from the go-mail/mail project.
// https://github.com/go-mail/mail (MIT License)

package mail

import (
	"bytes"
	"fmt"
	"net/smtp"
	"slices"
)

// loginAuth is an smtp.Auth that implements the LOGIN authentication mechanism.
type loginAuth struct {
	username string
	password string
	host     string
}

func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	if !server.TLS {
		if !slices.Contains(server.Auth, "LOGIN") {
			return "", nil, fmt.Errorf("gomail: unencrypted connection")
		}
	}
	if server.Name != a.host {
		return "", nil, fmt.Errorf("gomail: wrong host name")
	}
	return "LOGIN", nil, nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if !more {
		return nil, nil
	}

	switch {
	case bytes.Equal(fromServer, []byte("Username:")):
		return []byte(a.username), nil
	case bytes.Equal(fromServer, []byte("Password:")):
		return []byte(a.password), nil
	default:
		return nil, fmt.Errorf("gomail: unexpected server challenge: %s", fromServer)
	}
}
