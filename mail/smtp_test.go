package mail_test

// This fake SMTP server is inspired by https://github.com/andrewarrow/jungle_smtp
// and most of its functionality have been taken from the original repo and updated to
// work better for buffalo.

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

// fakeSMTPServer is our fake server that will be listening for SMTP connections.
type fakeSMTPServer struct {
	Listener net.Listener
	messages []string
	mutex    sync.Mutex
}

// Start listens for connections on the given port
func (s *fakeSMTPServer) Start(port string) error {
	for {
		conn, err := s.Listener.Accept()
		if err != nil {
			return err
		}

		s.Handle(&fakeSMTPConnection{
			conn:    conn,
			address: conn.RemoteAddr().String(),
			time:    time.Now().Unix(),
			bufin:   bufio.NewReader(conn),
			bufout:  bufio.NewWriter(conn),
		})
	}
}

// Handle a connection from a client
func (s *fakeSMTPServer) Handle(c *fakeSMTPConnection) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.messages = append(s.messages, "")

	s.readHello(c)
	s.readSender(c)
	s.readRecipients(c)
	s.readData(c)

	c.conn.Close()
}

// Requests and notifies readed the Hello
func (s *fakeSMTPServer) readHello(c *fakeSMTPConnection) {
	c.write("220 Welcome")
	text := c.read()
	s.addMessageLine(text)

	c.write("250 Received")
}

// readSender reads the Sender from the connection
func (s *fakeSMTPServer) readSender(c *fakeSMTPConnection) {
	text := c.read()
	s.addMessageLine(text)
	c.write("250 Sender")
}

// readRecipients reads recipients from the connection
func (s *fakeSMTPServer) readRecipients(c *fakeSMTPConnection) {
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

// readData reads the message data.
func (s *fakeSMTPServer) readData(c *fakeSMTPConnection) {
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

// addMessageLine ads a line to the last message
func (s *fakeSMTPServer) addMessageLine(text string) {
	s.messages[len(s.Messages())-1] = s.LastMessage() + text
}

// LastMessage returns the last message on the server
func (s *fakeSMTPServer) LastMessage() string {
	if len(s.Messages()) == 0 {
		return ""
	}

	return s.Messages()[len(s.Messages())-1]
}

// Messages returns the list of messages on the server
func (s *fakeSMTPServer) Messages() []string {
	return s.messages
}

// Clear the server messages
func (s *fakeSMTPServer) Clear() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.messages = []string{}
}

// newFakeSMTPServer returns a pointer to a new fakeSMTPServer instance listening on the given port.
func newFakeSMTPServer(port string) (*fakeSMTPServer, error) {
	s := &fakeSMTPServer{messages: []string{}}

	listener, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		return s, err
	}
	s.Listener = listener
	return s, nil
}

// fakeSMTPConnection of a client with our server
type fakeSMTPConnection struct {
	conn    net.Conn
	address string
	time    int64
	bufin   *bufio.Reader
	bufout  *bufio.Writer
}

// write something to the client on the connection
func (c *fakeSMTPConnection) write(s string) {
	c.bufout.WriteString(s + "\r\n")
	c.bufout.Flush()
}

// read a string from the connected client
func (c *fakeSMTPConnection) read() string {
	reply, err := c.bufin.ReadString('\n')

	if err != nil {
		fmt.Println("e ", err)
	}
	return reply
}
