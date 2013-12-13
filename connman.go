package main

import (
	"fmt"
	"net"
	"time"
)

// ConnMan is a connection manager
type ConnMan struct {
	ConnChan chan net.Conn
}

func NewConnMan() *ConnMan {
	c := &ConnMan{
		ConnChan: make(chan net.Conn),
	}

	go c.Start()
	return c
}

// Start handling connections
func (c *ConnMan) Start() {
	for conn := range c.ConnChan {
		go c.handleConnection(conn)
	}
}

func (c *ConnMan) handleConnection(conn net.Conn) {
	// Check the first byte for the magic byte
	b := make([]byte, 1)
	_, err := conn.Read(b)
	if err != nil {
		// I'm not sure what we'd do in this case...
		conn.Close()
	}

	if b[0] != MagicByte {
		conn.Write([]byte{ERR_MAGIC_BYTE})

		// There's not much we can do if
		// the magic byte is wrong.
		conn.Close()
	}

	conn.Write([]byte("Okay, I think that worked."))
	time.Sleep(time.Second * 2)
	b = make([]byte, 100)
	conn.Read(b)
	fmt.Println(string(b))
}
