package main

import (
	"net"
)

func (i *instance) handleGet(conn net.Conn) {
	buf := make([]byte, 1)

	conn.Read(buf)
	keyLength := buf[0]

	key := make([]byte, keyLength)
	conn.Read(key)

	value, _, err := i.db.Get(string(key))

	if err != nil {
		conn.Write([]byte{0})
	} else {
		conn.Write([]byte{byte(len(value))})
		conn.Write([]byte(value))
	}
}
