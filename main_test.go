package main

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"testing"
	"time"
)

func write(conn net.Conn, key string, value string) {
	fmt.Fprintf(conn, GenerateCommand(OP_SET, key, value))
}

func read(conn net.Conn, key string) string {
	fmt.Fprintf(conn, GenerateCommand(OP_GET, key))
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
	fmt.Fprintf(conn, GenerateCommand(OP_CLEARRANGE, "\x00", "\xff"))
}

func Test1(t *testing.T) {
	i := NewInstance(":12345", "")
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

func TestCommandLog(t *testing.T) {
	f, err := ioutil.TempFile(os.TempDir(), "fickle")
	if err != nil {
		t.Fatal("Couldn't open temp file")
	}

	i := NewInstance(":12346", f.Name())
	go i.Start()

	time.Sleep(time.Millisecond * 100) // Wait for it to start up
	conn, err := net.Dial("tcp", ":12346")
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

	conn.Close()

	i = NewInstance(":12347", f.Name())
	go i.Start()

	time.Sleep(time.Millisecond * 100) // Wait for it to start up
	conn, err = net.Dial("tcp", ":12347")
	if err != nil {
		t.Error(err)
	}
	if r := read(conn, "foo"); r != "bar" {
		t.Errorf("Bad read! Got %v", r)
	}

}

func TestReplica(t *testing.T) {
	i := NewInstance(":12348", "")
	j := NewInstance(":12349", "")

	go j.Start()
	go i.Start()

	time.Sleep(time.Millisecond * 100) // Wait for the replica to start

	i.AddReplica(":12349")

	// Sending to the primary
	conn, err := net.Dial("tcp", ":12348")
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
	conn.Close()

	// Reading from the replica
	conn, err = net.Dial("tcp", ":12349")
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
	conn.Close()
}
