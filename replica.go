package main

import (
	"fmt"
	"net"
)

type Replica struct {
	conn *net.Conn
}

// Sends a command to a replica and returns the error
// code from the replica
func (r *Replica) Send(command string) byte {
	fmt.Fprint((*r.conn), command)
	b := make([]byte, 1)
	(*r.conn).Read(b)
	return b[0]
}
