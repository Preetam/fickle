package main

import (
	"log"
	"net"
	"testing"
	"time"
)

func init() {
	i := new(instance)
	j := new(instance)
	go i.Start(":8080")
	go j.Start(":8081")
	time.Sleep(time.Second)
	err := i.AddReplica(":8081")
	if err != nil {
		log.Println(err.Error())
	}
}

func TestSimpleSetGetDelete(t *testing.T) {
	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		t.Error(err)
	}

	_, err = conn.Write([]byte("s\x01\x03\x06foofoobar"))
	if err != nil {
		t.Error(err)
	}

	_, err = conn.Write([]byte("g\x01\x03foo"))
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

	_, err = conn.Write([]byte("c\x01\x03foo"))
	if err != nil {
		t.Error(err)
	}

	_, err = conn.Write([]byte("g\x01\x03foo"))
	if err != nil {
		t.Error(err)
	}

	conn.Read(lenBuf)
	value = make([]byte, lenBuf[0])
	conn.Read(value)

	if string(value) != "" {
		t.Errorf(`Expected value to be "" got "%v"`, string(value))
	}

	conn.Close()
}

func TestReplica(t *testing.T) {
	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		t.Error(err)
	}

	_, err = conn.Write([]byte("s\x01\x03\x06foofoobar"))
	if err != nil {
		t.Error(err)
	}
	time.Sleep(time.Second)
	conn, err = net.Dial("tcp", ":8081")

	_, err = conn.Write([]byte("g\x01\x03foo"))
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

func TestReplicaSimpleSetGetDelete(t *testing.T) {
	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		t.Error(err)
	}

	_, err = conn.Write([]byte("s\x01\x03\x06foofoobar"))
	if err != nil {
		t.Error(err)
	}

	time.Sleep(time.Second)
	conn, err = net.Dial("tcp", ":8081")

	_, err = conn.Write([]byte("g\x01\x03foo"))
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

	conn, err = net.Dial("tcp", ":8080")

	_, err = conn.Write([]byte("c\x01\x03foo"))
	if err != nil {
		t.Error(err)
	}

	time.Sleep(time.Second)
	conn, err = net.Dial("tcp", ":8081")

	_, err = conn.Write([]byte("g\x01\x03foo"))
	if err != nil {
		t.Error(err)
	}

	conn.Read(lenBuf)
	value = make([]byte, lenBuf[0])
	conn.Read(value)

	if string(value) != "" {
		t.Errorf(`Expected value to be "" got "%v"`, string(value))
	}

	conn.Close()
}
