package main

const MagicByte = 0x14

const (
	ERR_NO_ERROR byte = iota
	ERR_MAGIC_BYTE
)

type Operation byte

const (
	OP_GET Operation = iota
	OP_SET
	OP_CLEAR
	OP_GETRANGE
	OP_CLEARRANGE
)
