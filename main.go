package main

import (
	"flag"
	"log"
	"net"
	"runtime"

	"github.com/PreetamJinka/lexicon"
)

type instance struct {
	db         *lexicon.Lexicon
	replicas   map[string]net.Conn
	listenAddr string
}

func (i *instance) Start(addr string) {
	i.listenAddr = addr
	i.db = lexicon.New()
	i.replicas = make(map[string]net.Conn)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		// handle error
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
			continue
		}
		go i.handleConnection(conn)
	}
}

func (i *instance) AddReplica(addr string) error {
	if _, present := i.replicas[addr]; !present {
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			return err
		} else {
			i.replicas[addr] = conn
		}
	}

	return nil
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	listenAddr := flag.String("listen", ":8080", "TCP address to listen on")
	debugHTTP := flag.Bool("debug-http", false, "Start an HTTP server for debugging")
	flag.Parse()

	i := new(instance)
	if *debugHTTP {
		StartHttpDebug()
	}
	i.Start(*listenAddr)
}

func (i *instance) handleConnection(conn net.Conn) {
	for {
		// Check the first character
		buf := make([]byte, 1)

		// Read the first character
		_, err := conn.Read(buf)
		if err != nil {
			log.Println("Error reading: ", err.Error())
			return
		}

		switch buf[0] {
		case 's':
			i.handleSet(conn)

		case 'g':
			i.handleGet(conn)

		case 'c':
			i.handleClear(conn)
	}
}
