package main

const MagicByte = 0x14
const HEADER_LEN = 6

const (
	ERR_NO_ERROR byte = iota
	ERR_MAGIC_BYTE
	ERR_BAD_HEADER
	ERR_INVALID_OP
	ERR_BAD_BODY
	ERR_INTERNAL
)

type Operation byte

const (
	OP_GET Operation = iota
	OP_SET
	OP_CLEAR
	OP_GETRANGE
	OP_CLEARRANGE
	OP_MAX_VALID
)
