package compiler

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"vmlite/code"
)

var stack = make([]interface{}, 2048)
var sp int = 0

func PrintByteCode(bc []code.Opcode, co_consts []interface{}) string {
	var out bytes.Buffer
	count := 0
	var ip int = -1

	for ip < len(bc)-1 {
		ip += 1
		c := bc[ip]
		count += 1
		switch c {
		case code.PUSH:
			i := binary.BigEndian.Uint32(bc[ip+1:])
			ip += 4
			v := co_consts[i]
			out.WriteString(fmt.Sprintf("%d\t%v\t%v\n", count, code.CodeMap[c], v))
		case code.ADD:
			out.WriteString(fmt.Sprintf("%d\t%v\n", count, code.CodeMap[c]))
		case code.SUB:
			out.WriteString(fmt.Sprintf("%d\t%v\n", count, code.CodeMap[c]))
		case code.MUL:
			out.WriteString(fmt.Sprintf("%d\t%v\n", count, code.CodeMap[c]))
		case code.DIV:
			out.WriteString(fmt.Sprintf("%d\t%v\n", count, code.CodeMap[c]))
		}
	}
	return out.String()
}

func push(v interface{}) {
	stack[sp] = v
	sp += 1
}

func pop() interface{} {
	v := stack[sp-1]
	sp -= 1
	return v
}
