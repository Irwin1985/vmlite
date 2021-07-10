package vm

import (
	"encoding/binary"
	"fmt"
	"strings"
	"vmlite/code"
)

const STACK_SIZE = 2048

type OpCodeFn func() error

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
	vm.mapCode[code.BOOL] = vm.OpBoolFn
	vm.mapCode[code.CMP] = vm.OpBinaryFn
	vm.mapCode[code.UNARY] = vm.OpUnaryFn
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
		err := opFn()
		if err != nil {
			return err
		}
	}
	return nil
}

func (vm *VM) OpPushFn() error {
	i := binary.BigEndian.Uint32(vm.co_codes[vm.ip:])
	vm.ip += 4
	v := vm.co_consts[i]
	vm.push(v)
	return nil
}

func (vm *VM) OpBoolFn() error {
	op := vm.co_codes[vm.ip]
	vm.ip += 1 // advance the instruction pointer
	if op == code.TRUE {
		vm.push(true)
	} else {
		vm.push(false)
	}
	return nil
}

func (vm *VM) OpBinaryFn() error {

	// get the comparison operator
	op := vm.co_codes[vm.ip]
	vm.ip += 1 // dont forget advance the instruction pointer

	switch op {
	case code.ADD:
		r := vm.pop()
		l := vm.pop()
		tr := typeOf(r)
		tl := typeOf(l)
		if tr == 'c' && tl == 'c' {
			vm.push(l.(string) + r.(string))
			return nil
		}
		if tr == 'n' && tl == 'n' {
			vm.push(l.(float64) + r.(float64))
			return nil
		}
		return fmt.Errorf("operands must be two numbers or two strings")
	case code.SUB:
		r := vm.pop()
		l := vm.pop()
		tr := typeOf(r)
		tl := typeOf(l)
		if tr == 'c' && tl == 'c' {
			vm.push(strings.TrimRight(l.(string), " ") + r.(string))
			return nil
		}
		if tr == 'n' && tl == 'n' {
			vm.push(l.(float64) - r.(float64))
			return nil
		}
		return fmt.Errorf("operands must be two numbers or two strings")
	case code.MUL:
		r, l := vm.popFloat()
		vm.push(l * r)
	case code.DIV:
		r, l := vm.popFloat()
		if r == 0 {
			return fmt.Errorf("division by zero")
		}
		vm.push(l / r)
	case code.LT:
		r, l := vm.popFloat()
		vm.push(l < r)
	case code.LEQ:
		r, l := vm.popFloat()
		vm.push(l <= r)
	case code.GT:
		r, l := vm.popFloat()
		vm.push(l > r)
	case code.GEQ:
		r, l := vm.popFloat()
		vm.push(l >= r)
	case code.EQ:
		r, l := vm.popFloat()
		vm.push(l == r)
	case code.NEQ:
		r, l := vm.popFloat()
		vm.push(l != r)
	case code.AND:
		r, l := vm.popBoolean()
		vm.push(l && r)
	case code.OR:
		r, l := vm.popBoolean()
		vm.push(l || r)
	}
	return nil
}

func (vm *VM) OpStoreFn() error {
	i := binary.BigEndian.Uint32(vm.co_codes[vm.ip:])
	vm.ip += 4
	v := vm.pop()
	vm.co_values[i] = v
	return nil
}

func (vm *VM) OpLoadFn() error {
	i := binary.BigEndian.Uint32(vm.co_codes[vm.ip:])
	vm.ip += 4
	v := vm.co_values[i]
	vm.push(v)
	return nil
}

func (vm *VM) OpUnaryFn() error {
	op := vm.co_codes[vm.ip]
	vm.ip += 1 // advance the ip

	if op == code.UNEG {
		v := vm.pop().(float64)
		vm.push(-v)
	} else {
		v := vm.pop().(bool)
		if v {
			vm.push(false)
		} else {
			vm.push(true)
		}
	}
	return nil
}

func (vm *VM) OpPrintFn() error {
	fmt.Printf("%v\n", vm.pop())
	return nil
}

// VIRTUAL MACHINE HELPER FUNCTIONS
func (vm *VM) popFloat() (float64, float64) {
	r := vm.pop().(float64)
	l := vm.pop().(float64)
	return r, l
}

func (vm *VM) popBoolean() (bool, bool) {
	r := vm.pop().(bool)
	l := vm.pop().(bool)
	return r, l
}

func typeOf(v interface{}) byte {
	switch v.(type) {
	case float64:
		return 'n'
	case string:
		return 'c'
	case bool:
		return 'l'
	case nil:
		return 'x'
	default:
		return 'u'
	}
}
