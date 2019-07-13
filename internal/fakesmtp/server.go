package fakesmtp

// This server is inspired by https://github.com/andrewarrow/jungle_smtp
// and most of its functionality have been taken from the original repo and updated to
// work better for buffalo.

import (
	"bufio"
	"net"
	"strings"
	"sync"
	"time"
)

//Server is our fake server that will be listening for SMTP connections.
type Server struct {
	Listener net.Listener
	messages []string
	mutex    sync.Mutex
}

//Start listens for connections on the given port
func (s *Server) Start(port string) error {
	for {
		conn, err := s.Listener.Accept()
		if err != nil {
			return err
		}

		s.Handle(&Connection{
			conn:    conn,
			address: conn.RemoteAddr().String(),
			time:    time.Now().Unix(),
			bufin:   bufio.NewReader(conn),
			bufout:  bufio.NewWriter(conn),
		})
	}
}

//Handle a connection from a client
func (s *Server) Handle(c *Connection) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.messages = append(s.messages, "")

	s.readHello(c)
	s.readSender(c)
	s.readRecipients(c)
	s.readData(c)

	c.conn.Close()
}

//Requests and notifies readed the Hello
func (s *Server) readHello(c *Connection) {
	c.write("220 Welcome")
	text := c.read()
	s.addMessageLine(text)

	c.write("250 Received")
}

//readSender reads the Sender from the connection
func (s *Server) readSender(c *Connection) {
	text := c.read()
	s.addMessageLine(text)
	c.write("250 Sender")
}

//readRecipients reads recipients from the connection
func (s *Server) readRecipients(c *Connection) {
	text := c.read()
	s.addMessageLine(text)

	c.write("250 Recipient")
	text = c.read()
	for strings.Contains(text, "RCPT") {
		s.addMessageLine(text)
		c.write("250 Recipient")
		text = c.read()
	}
}

//readData reads the message data.
func (s *Server) readData(c *Connection) {
	c.write("354 Ok Send data ending with <CRLF>.<CRLF>")

	for {
		text := c.read()
		bytes := []byte(text)
		s.addMessageLine(text)
		// 46 13 10
		if bytes[0] == 46 && bytes[1] == 13 && bytes[2] == 10 {
			break
		}
	}
	c.write("250 server has transmitted the message")
}

//addMessageLine ads a line to the last message
func (s *Server) addMessageLine(text string) {
	s.messages[len(s.Messages())-1] = s.LastMessage() + text
}

//LastMessage returns the last message on the server
func (s *Server) LastMessage() string {
	if len(s.Messages()) == 0 {
		return ""
	}

	return s.Messages()[len(s.Messages())-1]
}

//Messages returns the list of messages on the server
func (s *Server) Messages() []string {
	return s.messages
}

//Clear the server messages
func (s *Server) Clear() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.messages = []string{}
}

//New returns a pointer to a new Server instance listening on the given port.
func New(port string) (*Server, error) {
	s := &Server{messages: []string{}}

	listener, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		return s, err
	}
	s.Listener = listener
	return s, nil
}
