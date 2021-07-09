package vm

import (
	"encoding/binary"
	"fmt"
	"vmlite/code"
)

const STACK_SIZE = 2048

type VM struct {
	co_codes  []code.Opcode
	co_consts []interface{}
	stack     []interface{}
	sp        int
}

func NewVM(co_codes []code.Opcode, co_consts []interface{}) *VM {
	vm := &VM{
		co_codes:  co_codes,
		co_consts: co_consts,
		stack:     make([]interface{}, STACK_SIZE),
		sp:        0,
	}
	return vm
}

func (vm *VM) push(v interface{}) {
	vm.stack[vm.sp] = v
	vm.sp += 1
}

func (vm *VM) pop() interface{} {
	v := vm.stack[vm.sp-1]
	vm.sp -= 1
	return v
}

func (vm *VM) TOS() interface{} {
	return vm.stack[vm.sp-1]
}

func (vm *VM) Run() error {
	var ip = -1

	for ip < len(vm.co_codes)-1 {
		ip += 1
		op := vm.co_codes[ip]
		switch op {
		case code.PUSH:
			i := binary.BigEndian.Uint32(vm.co_codes[ip+1:])
			ip += 4
			v := vm.co_consts[i]
			vm.push(v)
		case code.ADD, code.SUB, code.MUL, code.DIV:
			r := vm.pop().(float64)
			l := vm.pop().(float64)
			switch op {
			case code.ADD:
				vm.push(l + r)
			case code.SUB:
				vm.push(l - r)
			case code.MUL:
				vm.push(l * r)
			case code.DIV:
				if r == 0 {
					return fmt.Errorf("division by zero")
				}
				vm.push(l / r)
			}
		}
	}

	return nil
}
