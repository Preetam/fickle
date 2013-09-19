package main

import (
	"fmt"
	"log"
	"net"

	"github.com/PreetamJinka/lexicon"
)

var db = lexicon.New()

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
		go func() { i.handleConnection(conn) }()
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
		}
	}
}

func (i *instance) handleSet(conn net.Conn) {
	buf := make([]byte, 1)

	log.Println("SET")

	conn.Read(buf)
	keyLength := buf[0]
	log.Printf("keyLength: %d\n", keyLength)

	conn.Read(buf)
	valueLength := buf[0]
	log.Printf("valueLength: %d\n", valueLength)

	key := make([]byte, keyLength)
	value := make([]byte, valueLength)

	conn.Read(key)
	conn.Read(value)

	log.Printf("Setting %s => %s\n", key, value)
	i.db.Set(string(key), string(value))

	for rep, conn := range i.replicas {
		command := replicaSetCommandHelper(key, value)
		log.Printf("Sending to replica %s: %s\n", rep, command)
		conn.Write(command)
	}
}

func (i *instance) handleGet(conn net.Conn) {
	buf := make([]byte, 1)

	log.Println("GET")

	conn.Read(buf)
	keyLength := buf[0]
	log.Printf("keyLength: %d\n", keyLength)

	key := make([]byte, keyLength)
	conn.Read(key)

	value, _, err := i.db.Get(string(key))

	if err != nil {
		conn.Write([]byte{0})
	} else {
		conn.Write([]byte{byte(len(value))})
		conn.Write([]byte(value))
	}
}

func (i *instance) handleDelete(conn net.Conn) {
	buf := make([]byte, 1)

	log.Println("DELETE")

	conn.Read(buf)
	keyLength := buf[0]
	log.Printf("keyLength: %d\n", keyLength)

	key := make([]byte, keyLength)
	conn.Read(key)

	i.db.Remove(string(key))
}

func replicaSetCommandHelper(key, value []byte) []byte {
	buf := []byte{'s', byte(len(key)), byte(len(value))}
	return []byte(fmt.Sprintf("%s%s%s", buf, key, value))
}
