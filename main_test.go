package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"testing"
	"time"
)

func commandGenerator(op Operation, vars ...string) string {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, byte(MagicByte))
	binary.Write(buf, binary.LittleEndian, byte(op))
	for _, v := range vars {
		binary.Write(buf, binary.LittleEndian, uint16(len(v)))
	}
	// Pad until we have the minimum header length
	for buf.Len() < HEADER_LEN {
		buf.WriteByte(0x0)
	}
	for _, v := range vars {
		buf.WriteString(v)
	}
	return buf.String()
}

func write(conn net.Conn, key string, value string) {
	fmt.Fprintf(conn, commandGenerator(OP_SET, key, value))
}

func read(conn net.Conn, key string) string {
	fmt.Fprintf(conn, commandGenerator(OP_GET, key))
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
	return b[0] == ERR_NO_ERROR
}

func clearAll(conn net.Conn) {
	fmt.Fprintf(conn, commandGenerator(OP_CLEARRANGE, "\x00", "\xff"))
}

func Test1(t *testing.T) {
	i := NewInstance(":12345")
	go i.Start()
	time.Sleep(time.Millisecond * 100) // Wait for it to start up
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
