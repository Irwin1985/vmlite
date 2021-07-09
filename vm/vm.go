package vm

import (
	"encoding/binary"
	"fmt"
	"vmlite/code"
)

const STACK_SIZE = 2048

type OpCodeFn func(code.Opcode) error

type VM struct {
	co_codes  []code.Opcode
	co_consts []interface{}
	co_names  []string
	co_values []interface{}
	stack     []interface{}
	sp        int
	ip        int
	mapCode   map[code.Opcode]OpCodeFn
}

func NewVM(co_codes []code.Opcode, co_consts []interface{}, co_names []string, co_values []interface{}) *VM {
	vm := &VM{
		co_codes:  co_codes,
		co_consts: co_consts,
		co_names:  co_names,
		co_values: co_values,
		stack:     make([]interface{}, STACK_SIZE),
		sp:        0,
		mapCode:   make(map[byte]OpCodeFn),
	}
	// register semantic opcode
	vm.mapCode[code.PUSH] = vm.OpPushFn
	vm.mapCode[code.ADD] = vm.OpBinaryFn
	vm.mapCode[code.SUB] = vm.OpBinaryFn
	vm.mapCode[code.MUL] = vm.OpBinaryFn
	vm.mapCode[code.DIV] = vm.OpBinaryFn
	vm.mapCode[code.UNEG] = vm.OpUnaryFn
	vm.mapCode[code.STORE] = vm.OpStoreFn
	vm.mapCode[code.LOAD] = vm.OpLoadFn
	vm.mapCode[code.PRINT] = vm.OpPrintFn
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
	if vm.sp > 0 {
		return vm.stack[vm.sp-1]
	}
	return nil
}

func (vm *VM) Run() error {

	for vm.ip < len(vm.co_codes) {
		op := vm.co_codes[vm.ip]
		vm.ip += 1
		opFn := vm.mapCode[op]
		if opFn == nil {
			return fmt.Errorf("unknown opcode: <%v, %v>", op, code.CodeMap[op])
		}
		err := opFn(op)
		if err != nil {
			return err
		}
	}
	return nil
}

func (vm *VM) OpPushFn(op code.Opcode) error {
	i := binary.BigEndian.Uint32(vm.co_codes[vm.ip:])
	vm.ip += 4
	v := vm.co_consts[i]
	vm.push(v)
	return nil
}

func (vm *VM) OpBinaryFn(op code.Opcode) error {
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
	return nil
}

func (vm *VM) OpStoreFn(op code.Opcode) error {
	i := binary.BigEndian.Uint32(vm.co_codes[vm.ip:])
	vm.ip += 4
	v := vm.pop()
	vm.co_values[i] = v
	return nil
}

func (vm *VM) OpLoadFn(op code.Opcode) error {
	i := binary.BigEndian.Uint32(vm.co_codes[vm.ip:])
	vm.ip += 4
	v := vm.co_values[i]
	vm.push(v)
	return nil
}

func (vm *VM) OpUnaryFn(op code.Opcode) error {
	v := vm.pop().(float64)
	vm.push(-v)
	return nil
}

func (vm *VM) OpPrintFn(op code.Opcode) error {
	fmt.Printf("%v\n", vm.pop())
	return nil
}
