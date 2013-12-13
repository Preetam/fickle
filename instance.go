package main

import (
	"net"

	"github.com/PreetamJinka/lexicon"
)

// Instance is a fickle instance
type Instance struct {
	db         *lexicon.Lexicon
	replicas   map[string]net.Conn
	listenAddr string
	connman    *ConnMan
}

func NewInstance(addr string) *Instance {
	return &Instance{
		db:         lexicon.New(),
		replicas:   make(map[string]net.Conn),
		listenAddr: addr,
		connman:    NewConnMan(),
	}
}

func (i *Instance) Start() {
	ln, err := net.Listen("tcp", i.listenAddr)
	if err != nil {
		panic("Couldn't listen on " + i.listenAddr)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			// ignore it
			continue
		}

		// Pass the connection on to the
		// connection manager.
		i.connman.ConnChan <- conn
	}
}
