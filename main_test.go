package main

import (
	"net"
	"testing"
	"time"
)

func init() {
	go main()
	time.Sleep(time.Second)
}

func TestSimpleSetGet(t *testing.T) {
	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		t.Error(err)
	}

	_, err = conn.Write([]byte("s\x03\x06foofoobar"))
	if err != nil {
		t.Error(err)
	}

	_, err = conn.Write([]byte("g\x03foo"))
	if err != nil {
		t.Error(err)
	}

	lenBuf := make([]byte, 1)
	conn.Read(lenBuf)

	value := make([]byte, lenBuf[0])
	conn.Read(value)

	if string(value) != "foobar" {
		t.Errorf(`Expected value to be "foobar" got "%v"`, string(value))
	}

	conn.Close()
}
