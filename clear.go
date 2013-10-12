package main

import (
	"net"
)

func (i *instance) handleClear(conn net.Conn) {

	// repcmd is basically the same
	// command we're getting.
	repcmd := []byte("c")
	buf := make([]byte, 1)

	conn.Read(buf)
	repcmd = append(repcmd, buf...)

	for keyCount := int(buf[0]); keyCount > 0; keyCount-- {
		conn.Read(buf)
		repcmd = append(repcmd, buf...)

		keyLength := buf[0]

		key := make([]byte, keyLength)

		conn.Read(key)
		repcmd = append(repcmd, key...)

		i.db.Remove(ComparableString(key))

		for _, repConn := range i.replicas {
			// echo the repcmd to every replica
			repConn.Write(repcmd)
		}

	}
}
