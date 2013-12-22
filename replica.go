/*

TODO

*/

package main

import (
	"net"
)

type Replica struct {
	conn *net.Conn
}

func (r *Replica) Send(command string) byte {
	return 0
}
