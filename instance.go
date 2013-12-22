package main

import (
	"encoding/binary"
	"io"
	"log"
	"net"
	"os"

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
	replicas   map[string]*Replica
	listenAddr string
	connman    *ConnMan
	commandLog string
	file       *os.File
}

func NewInstance(addr string, log string) *Instance {
	i := &Instance{
		db:         lexicon.New(compareStrings),
		replicas:   make(map[string]*Replica),
		listenAddr: addr,
	}
	i.connman = NewConnMan(i)

	if log != "" {
		f, err := os.OpenFile(log, os.O_RDWR, 0777)
		if err != nil {

		}

		i.file = f
	}

	return i
}

func (i *Instance) Start() {
	defer func() {
		if r := recover(); r != nil {
			i.file.Close()
			log.Println("Recovered from a panic:", r)
		}
	}()
	ln, err := net.Listen("tcp", i.listenAddr)
	if err != nil {
		panic("Couldn't listen on " + i.listenAddr)
	}

	// Reload from the log
	if i.file != nil {
		conn, err := net.Dial("tcp", i.listenAddr)
		if err == nil {
			_, err := io.Copy(conn, i.file)
			if err != nil {
				log.Println("Error reading from command log:", err)
			} else {
				i.commandLog = i.file.Name()
			}
		}
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

func (i *Instance) AddReplica(address string) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		i.replicas[address].conn = &conn
	}
}

// This is a pretty lame way to do it,
// but I'll fix it later :)
func (i *Instance) Execute(c Command) (byte, []byte) {

	switch c.Type {
	case OP_GET:
		return i.Get(c.Var1)
	case OP_SET:
		if err := i.LogCommand(c); err != nil {
			return ERR_INTERNAL, nil
		}
		i.Replicate(c)
		return i.Set(c.Var1, c.Var2)
	case OP_CLEAR:
		if err := i.LogCommand(c); err != nil {
			return ERR_INTERNAL, nil
		}
		i.Replicate(c)
		return i.Clear(c.Var1)
	case OP_GETRANGE:
		return i.GetRange(c.Var1, c.Var2)
	case OP_CLEARRANGE:
		if err := i.LogCommand(c); err != nil {
			return ERR_INTERNAL, nil
		}
		i.Replicate(c)
		return i.ClearRange(c.Var1, c.Var2)
	}

	return ERR_INVALID_OP, nil
}

func (i *Instance) Replicate(c Command) {
	cmdStr := ""
	if c.Var2 != "" {
		cmdStr = GenerateCommand(c.Type, c.Var1, c.Var2)
	} else {
		cmdStr = GenerateCommand(c.Type, c.Var1)
	}

	for _, replica := range i.replicas {
		replica.Send(cmdStr)
	}
}

func (i *Instance) LogCommand(c Command) error {
	if i.commandLog == "" {
		return nil
	}

	cmdStr := ""
	if c.Var2 != "" {
		cmdStr = GenerateCommand(c.Type, c.Var1, c.Var2)
	} else {
		cmdStr = GenerateCommand(c.Type, c.Var1)
	}

	_, err := i.file.WriteString(cmdStr)
	if err != nil {
		log.Println(err)
		return err
	}

	i.file.Sync()

	return nil
}

func (i *Instance) Get(key string) (resErr byte, resBody []byte) {
	r, err := i.db.Get(key)
	if err == lexicon.ErrKeyNotPresent {
		resErr = ERR_NO_ERROR
	}

	cs, ok := r.(string)
	if !ok {
		resErr = ERR_INTERNAL
		return
	}

	resBody = stringToByteArray(cs)

	return
}

func (i *Instance) Set(key, val string) (resErr byte, resBody []byte) {
	i.db.Set(key, val)

	resErr = ERR_NO_ERROR
	return
}

func (i *Instance) Clear(key string) (resErr byte, resBody []byte) {
	i.db.Remove(key)
	resErr = ERR_NO_ERROR
	return
}

func (i *Instance) GetRange(start, end string) (resErr byte, resBody []byte) {
	kv := i.db.GetRange(start, end)
	resErr = ERR_NO_ERROR
	resBody = keyValueArrayToByteArray(kv)

	return
}

func (i *Instance) ClearRange(start, end string) (resErr byte, resBody []byte) {
	i.db.ClearRange(start, end)
	resErr = ERR_NO_ERROR

	return
}

func stringToByteArray(cs string) []byte {
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
		out = append(out, stringToByteArray(i.Key.(string))...)
		out = append(out, stringToByteArray(i.Value.(string))...)
	}

	return out
}
