package vm

import (
	"fmt"

	"github.com/badc0re/hprog/chunk"
	"github.com/badc0re/hprog/codes"
	"github.com/badc0re/hprog/lexer"
	"github.com/badc0re/hprog/stack"
	"github.com/badc0re/hprog/token"
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
		case codes.INSTRUC_CONSTANT:
			// fmt.Println("CONST")
			constant := v.ReadConstant()
			// value.PrintValue(constant)
			v.stack.Push(constant)
			break
		case codes.INSTRUC_ADDITION:
			v.binaryOP("+")
		case codes.INSTRUC_SUBSTRACT:
			v.binaryOP("-")
		case codes.INSTRUC_MULTIPLY:
			v.binaryOP("*")
		case codes.INSTRUC_DIVIDE:
			v.binaryOP("/")
		case codes.INSTRUC_NEGATE:
			v.stack.Push(-v.stack.Pop())
		case codes.INSTRUC_RETURN:
			fmt.Println("RETURN")
			fmt.Println("STACK POP:", v.stack.Pop())
			return INTER_OK
		}
	}
}

/*
func (v *VM) Interpret(chk *chunk.Chunk) INTER_RESULT {
	v.chunk = chk
	v.counter = 0
	v.ip = &v.chunk.Code[v.counter]
	return v.run()
}
*/
func Compile(source string, chk *chunk.Chunk) INTER_RESULT {
	lex := lexer.Init(source)
	//p := Parser{}

	for tkn := range lex.Consume() {

		switch tkn.Type {
		case token.NUMBER:
			// chk.AddVariable(tkn.Value)
		case token.ERR:
			return INTER_COMPILE_ERROR
		}
		if tkn.Type == token.ERR {
		}
	}
	return INTER_OK
}

type Parser struct {
	current  *token.Token
	previous *token.Token
	perror   bool
	ppanic   bool
}

func (v *VM) Interpret(source string) INTER_RESULT {
	chk := chunk.Chunk{}

	if Compile(source, &chk) == INTER_COMPILE_ERROR {
		// parser.ppanic = true
		// parser.perror = true
		return INTER_COMPILE_ERROR
	}
	if len(chk.Code) != 0 {
		v.chunk = &chk
		v.counter = 0
		v.ip = &v.chunk.Code[v.counter]
		return v.run()
	}
	return INTER_OK
}
