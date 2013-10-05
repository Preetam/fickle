package main

import (
	"encoding/json"
	"log"
	"net"
	"runtime"

	"github.com/PreetamJinka/lexicon"
)

var m runtime.MemStats

type instance struct {
	db       *lexicon.Lexicon
	replicas map[string]net.Conn
}

func (i *instance) Start(addr string) {
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
	new(instance).Start(":8080")
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

		case 'd':
			i.handleDelete(conn)

		case '0':
			dumpStats()

		case '1':
			runtime.GC()
		}
	}
}

func dumpStats() {
	runtime.ReadMemStats(&m)
	marshalled, err := json.Marshal(&m)
	if err == nil {
		log.Println("\n\n", string(marshalled))
	}
}
