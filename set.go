package main

import (
	"fmt"
	"net"
)

func replicaSetCommandHelper(key, value []byte) []byte {
	buf := []byte{'s', byte(len(key)), byte(len(value))}
	return []byte(fmt.Sprintf("%s%s%s", buf, key, value))
}

func (i *instance) handleSet(conn net.Conn) {
	buf := make([]byte, 1)

	conn.Read(buf)
	keyLength := buf[0]

	conn.Read(buf)
	valueLength := buf[0]

	key := make([]byte, keyLength)
	value := make([]byte, valueLength)

	conn.Read(key)
	conn.Read(value)

	i.db.Set(ComparableString(key), ComparableString(value))

	for _, conn := range i.replicas {
		command := replicaSetCommandHelper(key, value)
		conn.Write(command)
	}
}
