package main

import (
	"encoding/binary"
	"net"

	"github.com/PreetamJinka/lexicon"
)

func compareStrings(a, b interface{}) (result int) {
	defer func() {
		if r := recover(); r != nil {
			// Log it?
		}
	}()

	aStr := a.(string)
	bStr := b.(string)

	if aStr > bStr {
		result = 1
	}

	if aStr < bStr {
		result = -1
	}

	return
}

// Instance is a fickle instance
type Instance struct {
	db         *lexicon.Lexicon
	replicas   map[string]net.Conn
	listenAddr string
	connman    *ConnMan
}

func NewInstance(addr string) *Instance {
	i := &Instance{
		db:         lexicon.New(compareStrings),
		replicas:   make(map[string]net.Conn),
		listenAddr: addr,
	}
	i.connman = NewConnMan(i)

	return i
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

// This is a pretty lame way to do it,
// but I'll fix it later :)
func (i *Instance) Execute(c Command) (byte, []byte) {
	switch c.Type {
	case OP_GET:
		return i.Get(c.Var1)
	case OP_SET:
		return i.Set(c.Var1, c.Var2)
	case OP_CLEAR:
		return i.Clear(c.Var1)
	case OP_GETRANGE:
		return i.GetRange(c.Var1, c.Var2)
	case OP_CLEARRANGE:
		return i.ClearRange(c.Var1, c.Var2)
	}

	return ERR_INVALID_OP, nil
}

func (i *Instance) Get(key ComparableString) (resErr byte, resBody []byte) {
	r, err := i.db.Get(key)
	if err == lexicon.ErrKeyNotPresent {
		resErr = ERR_NO_ERROR
	}

	cs, ok := r.(ComparableString)
	if !ok {
		resErr = ERR_INTERNAL
		return
	}

	resBody = comparableStringToByteArray(cs)

	return
}

func (i *Instance) Set(key ComparableString, val ComparableString) (resErr byte, resBody []byte) {
	i.db.Set(key, val)

	resErr = ERR_NO_ERROR
	return
}

func (i *Instance) Clear(key ComparableString) (resErr byte, resBody []byte) {
	i.db.Remove(key)
	resErr = ERR_NO_ERROR
	return
}

func (i *Instance) GetRange(start, end ComparableString) (resErr byte, resBody []byte) {
	kv := i.db.GetRange(start, end)
	resErr = ERR_NO_ERROR
	resBody = keyValueArrayToByteArray(kv)

	return
}

func (i *Instance) ClearRange(start, end ComparableString) (resErr byte, resBody []byte) {
	i.db.ClearRange(start, end)
	resErr = ERR_NO_ERROR

	return
}

func comparableStringToByteArray(cs ComparableString) []byte {
	size := uint16(len(cs))
	sizeBuf := make([]byte, 2)

	binary.LittleEndian.PutUint16(sizeBuf, size)
	return append(sizeBuf, []byte(cs)...)
}

func keyValueArrayToByteArray(kv []lexicon.KeyValue) []byte {
	size := uint64(len(kv))
	out := make([]byte, 8)
	binary.LittleEndian.PutUint64(out, size)

	for _, i := range kv {
		out = append(out, comparableStringToByteArray(i.Key.(ComparableString))...)
		out = append(out, comparableStringToByteArray(i.Value.(ComparableString))...)
	}

	return out
}
