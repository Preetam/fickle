package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"testing"
	"time"
)

func write(conn net.Conn, key string, value string) {
	binary.Write(conn, binary.LittleEndian, byte(0x14))         // don't change this
	binary.Write(conn, binary.LittleEndian, byte(1))            // 1 => set, 2 => clear
	binary.Write(conn, binary.LittleEndian, uint16(len(key)))   // length of the key
	binary.Write(conn, binary.LittleEndian, uint16(len(value))) // length of the value
	fmt.Fprintf(conn, key)                                      // the key
	fmt.Fprintf(conn, value)                                    // the value
}

func read(conn net.Conn, key string) string {
	binary.Write(conn, binary.LittleEndian, byte(0x14))       // don't change this
	binary.Write(conn, binary.LittleEndian, byte(0))          // 1 => set, 2 => clear
	binary.Write(conn, binary.LittleEndian, uint16(len(key))) // length of the key
	binary.Write(conn, binary.LittleEndian, uint16(1))        // length of the value
	fmt.Fprintf(conn, key)                                    // the key
	fmt.Fprintf(conn, "0")
	if verify(conn) {
		var size uint16
		binary.Read(conn, binary.LittleEndian, &size)
		b := make([]byte, size)
		_, err := conn.Read(b)
		if err == nil {
			return string(b)
		}
	}
	return ""
}

func verify(conn net.Conn) bool {
	b := make([]byte, 1)
	conn.Read(b)
	return b[0] == 0
}

func clearAll(conn net.Conn) {
	binary.Write(conn, binary.LittleEndian, byte(0x14)) // don't change this
	binary.Write(conn, binary.LittleEndian, byte(4))    // 1 => set, 2 => clear
	binary.Write(conn, binary.LittleEndian, uint16(1))  // length of the key
	binary.Write(conn, binary.LittleEndian, uint16(1))  // length of the value
	fmt.Fprintf(conn, "\x00")                           // the key
	fmt.Fprintf(conn, "\xff")                           // the value
}

func Test1(t *testing.T) {
	i := NewInstance(":12345")
	go i.Start()
	time.Sleep(time.Millisecond * 100)
	conn, err := net.Dial("tcp", ":12345")
	if err != nil {
		t.Error(err)
	}
	write(conn, "foo", "bar")
	if !verify(conn) {
		t.Error("Bad write!")
	}
	if r := read(conn, "foo"); r != "bar" {
		t.Errorf("Bad read! Got %v", r)
	}
}
