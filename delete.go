package main

import (
	"fmt"
	"net"
)

func replicaDeleteCommandHelper(key []byte) []byte {
	buf := []byte{'d', byte(len(key))}
	return []byte(fmt.Sprintf("%s%s", buf, key))
}

func (i *instance) handleDelete(conn net.Conn) {
	buf := make([]byte, 1)

	conn.Read(buf)
	keyLength := buf[0]

	key := make([]byte, keyLength)
	conn.Read(key)

	i.db.Remove(ComparableString(key))

	for _, conn := range i.replicas {
		command := replicaDeleteCommandHelper(key)
		conn.Write(command)
	}
}
