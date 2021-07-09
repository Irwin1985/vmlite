package code

import "encoding/binary"

type Opcode = byte

const (
	PUSH Opcode = iota
	ADD
	SUB
	MUL
	DIV
	UNEG
	STORE
	LOAD
	PRINT
)

var CodeMap = map[Opcode]string{
	PUSH:  "PUSH",
	ADD:   "ADD",
	SUB:   "SUB",
	MUL:   "MUL",
	DIV:   "DIV",
	UNEG:  "UNEG",
	STORE: "STORE",
	LOAD:  "LOAD",
	PRINT: "PRINT",
}

func Make(op Opcode, args ...int) []byte {
	// 1 = the opcode + 4 bytes for int32 * len(args)
	// [0,  1, 2, 3, 4,   5, 6, 7, 8, 9, 10, 11, 12, ...]
	//  ^ | ^  ^  ^  ^  | ^  ^  ^  ^ |
	//  |       |			   |
	//  |       |			   |
	//  v       v			   v
	// opCode arg1			  arg2
	size := 1 + 4*len(args) // eg: 10 1 2 => 1 + 4 * 2 = 9
	b := make([]byte, size)
	b[0] = op
	for i, v := range args {
		// offset formula => i * 4 + 1
		o := i*4 + 1
		binary.BigEndian.PutUint32(b[o:], uint32(v))
	}

	return b
}
