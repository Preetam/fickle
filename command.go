package main

import (
	"bytes"
	"encoding/binary"
)

type Command struct {
	Type Operation
	Var1 string
	Var2 string
}

func RunCommand(i *Instance, op byte, var1 []byte, var2 []byte) []byte {
	c := Command{
		Type: Operation(op),
		Var1: string(var1),
		Var2: string(var2),
	}
	errByte, response := i.Execute(c)

	return append([]byte{errByte}, response...)
}

func GenerateCommand(op Operation, vars ...string) string {
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
