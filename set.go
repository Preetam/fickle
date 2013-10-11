package main

import (
	"net"
)

func (i *instance) handleSet(conn net.Conn) {
	repcmd := []byte("s")
	buf := make([]byte, 1)

	conn.Read(buf)
	repcmd = append(repcmd, buf...)

	for keyCount := int(buf[0]); keyCount > 0; keyCount-- {
		conn.Read(buf)
		repcmd = append(repcmd, buf...)

		keyLength := buf[0]

		conn.Read(buf)
		repcmd = append(repcmd, buf...)

		valueLength := buf[0]

		key := make([]byte, keyLength)
		value := make([]byte, valueLength)

		conn.Read(key)
		repcmd = append(repcmd, key...)

		conn.Read(value)
		repcmd = append(repcmd, value...)

		i.db.Set(ComparableString(key), ComparableString(value))
	}

	for _, repConn := range i.replicas {
		repConn.Write(repcmd)
	}
}
