package vm

import (
	"fmt"

	"github.com/badc0re/hprog/chunk"
	"github.com/badc0re/hprog/stack"
	"github.com/badc0re/hprog/value"
)

type VM struct {
	chunk   *chunk.Chunk
	ip      *interface{}
	counter int
	stack   stack.Stack
}

type INTER_RESULT int

const (
	INTER_ILLEGAL INTER_RESULT = iota

	INTER_OK
	INTER_COMPILE_ERROR
	INTER_RUNTIME_ERROR
)

func (vm *VM) InitVM() {

}

func (vm *VM) FreeVM() {

}

func (vm *VM) Move() interface{} {
	vm.ip = &vm.chunk.Code[vm.counter]
	vm.counter++
	return *vm.ip
}

func (vm *VM) ReadConstant() value.Value {
	vm.Move()
	index := (*vm.ip).(uint)
	return vm.chunk.Constants.Values[index]
}

func (vm *VM) binaryOP(op string) {
	b := vm.stack.Pop()
	a := vm.stack.Pop()
	switch op {
	case "+":
		vm.stack.Push(a + b)
	case "-":
		vm.stack.Push(a - b)
	case "/":
		vm.stack.Push(a / b)
	case "*":
		vm.stack.Push(a * b)
	}
}

func (v *VM) run() INTER_RESULT {
	for {
		instruct := v.Move()
		switch instruct {
		case chunk.OP_CONSTANT:
			fmt.Println("CONST")
			constant := v.ReadConstant()
			value.PrintValue(constant)
			v.stack.Push(constant)
			break
		case chunk.OP_ADD:
			v.binaryOP("+")
		case chunk.OP_SUBSTRACT:
			v.binaryOP("-")
		case chunk.OP_MULTIPLY:
			fmt.Println("AAAA")
			v.binaryOP("*")
		case chunk.OP_DIVIDE:
			v.binaryOP("/")
		case chunk.OP_NEGATE:
			v.stack.Push(-v.stack.Pop())
			fmt.Printf("%g\n", v.stack)
		case chunk.OP_RETURN:
			fmt.Println("RETURN")
			fmt.Println("STACK POP:", v.stack.Pop())
			return INTER_OK
		}
	}
}

func (v *VM) Interpret(chk *chunk.Chunk) INTER_RESULT {
	v.chunk = chk
	v.counter = 0
	v.ip = &v.chunk.Code[v.counter]
	return v.run()
}
