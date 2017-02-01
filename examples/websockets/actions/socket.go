package actions

import (
	"strings"
	"time"

	"github.com/gobuffalo/buffalo"
	"github.com/pkg/errors"
)

// Message stores incoming data from websocket connection buffer
type Message struct {
	Original  string    `json:"original"`
	Formatted string    `json:"formatted"`
	Received  time.Time `json:"received"`
}

// SocketHandler reads messages from the websocket connection buffer and returns
// the original message, formatted (uppercase) message, and received time in JSON
func SocketHandler(c buffalo.Context) error {
	conn, err := c.Websocket()
	if err != nil {
		return errors.WithStack(err)
	}
	for {
		// Read a message from the connection buffer.
		_, m, err := conn.ReadMessage()
		if err != nil {
			return errors.WithStack(err)
		}

		// Convert the bytes we received to a string.
		data := string(m)

		// Create a message and store the data.
		msg := Message{
			Original:  data,
			Formatted: strings.ToUpper(data),
			Received:  time.Now(),
		}

		// Encode the message to JSON and send it back.
		if err := conn.WriteJSON(msg); err != nil {
			return errors.WithStack(err)
		}
	}
}
