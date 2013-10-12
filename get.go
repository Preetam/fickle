package main

import (
	"net"
)

func (i *instance) handleGet(conn net.Conn) {

	// Read the key count
	buf := make([]byte, 1)
	conn.Read(buf)

	for keyCount := int(buf[0]); keyCount > 0; keyCount-- {
		conn.Read(buf)
		keyLength := buf[0]

		key := make([]byte, keyLength)
		conn.Read(key)

		value, err := i.db.Get(ComparableString(key))

		if err != nil {
			// err => key not found
			conn.Write([]byte{0})
		} else {
			conn.Write([]byte{byte(len(value.(ComparableString)))})
			conn.Write([]byte(value.(ComparableString)))
		}
	}
}
