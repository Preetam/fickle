package main

import (
	"log"
	"net"

	"github.com/PreetamJinka/lexicon"
)

var db = lexicon.New()

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		// handle error
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
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
			handleSet(conn)

		case 'g':
			handleGet(conn)

		case 'd':
			handleDelete(conn)
		}
	}
}

func handleSet(conn net.Conn) {
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
	db.Set(string(key), string(value))
}

func handleGet(conn net.Conn) {
	buf := make([]byte, 1)

	log.Println("GET")

	conn.Read(buf)
	keyLength := buf[0]
	log.Printf("keyLength: %d\n", keyLength)

	key := make([]byte, keyLength)
	conn.Read(key)

	value, _, err := db.Get(string(key))

	if err != nil {
		conn.Write([]byte{0})
	} else {
		conn.Write([]byte{byte(len(value))})
		conn.Write([]byte(value))
	}
}

func handleDelete(conn net.Conn) {
	buf := make([]byte, 1)

	log.Println("DELETE")

	conn.Read(buf)
	keyLength := buf[0]
	log.Printf("keyLength: %d\n", keyLength)

	key := make([]byte, keyLength)
	conn.Read(key)

	db.Remove(string(key))
}
