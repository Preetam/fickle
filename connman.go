package main

import (
	"bytes"
	"encoding/binary"
	"net"
)

type MessageHeader struct {
	Opcode      uint8
	Var1_length uint16
	Var2_length uint16
}

// ConnMan is a connection manager
type ConnMan struct {
	ConnChan chan net.Conn
	instance *Instance
}

func NewConnMan(i *Instance) *ConnMan {
	c := &ConnMan{
		ConnChan: make(chan net.Conn),
		instance: i,
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

	// If something goes wrong, we'll just return
	// and close the connection.
	defer conn.Close()
	for {
		// Check the first byte for the magic byte
		b := make([]byte, 1)
		_, err := conn.Read(b)
		if err != nil {
			// I'm not sure what we'd do in this case...
			return
		}

		if b[0] != MagicByte {
			conn.Write([]byte{ERR_MAGIC_BYTE})

			// There's not much we can do if
			// the magic byte is wrong.
			return
		}

		header := MessageHeader{}
		b = make([]byte, 5)
		_, err = conn.Read(b)
		if err != nil {
			// I'm not sure what we'd do in this case...
			return
		}
		r := bytes.NewReader(b)
		err = binary.Read(r, binary.LittleEndian, &header)
		if err != nil {
			conn.Write([]byte{ERR_BAD_HEADER})
			return
		}
		if header.Opcode >= byte(OP_MAX_VALID) {
			conn.Write([]byte{ERR_INVALID_OP})
			return
		}

		var1 := make([]byte, header.Var1_length)

		var2 := make([]byte, header.Var2_length)

		_, err = conn.Read(var1)
		if err != nil {
			conn.Write([]byte{ERR_BAD_BODY})
			return
		}
		if header.Var2_length > 0 {
			_, err = conn.Read(var2)
			if err != nil {
				conn.Write([]byte{ERR_BAD_BODY})
				return
			}
		}
		conn.Write(RunCommand(c.instance, header.Opcode, var1, var2))
	}
}
