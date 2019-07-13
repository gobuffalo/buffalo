package fakesmtp

import (
	"bufio"
	"fmt"
	"net"
)

//Connection of a client with our server
type Connection struct {
	conn    net.Conn
	address string
	time    int64
	bufin   *bufio.Reader
	bufout  *bufio.Writer
}

//write something to the client on the connection
func (c *Connection) write(s string) {
	c.bufout.WriteString(s + "\r\n")
	c.bufout.Flush()
}

//read a string from the connected client
func (c *Connection) read() string {
	reply, err := c.bufin.ReadString('\n')

	if err != nil {
		fmt.Println("e ", err)
	}
	return reply
}
