package vm

import (
	"fmt"

	"github.com/badc0re/hprog/chunk"
	"github.com/badc0re/hprog/codes"
	"github.com/badc0re/hprog/lexer"
	"github.com/badc0re/hprog/parser"
	"github.com/badc0re/hprog/stack"
	"github.com/badc0re/hprog/token"
	"github.com/badc0re/hprog/value"
)

var MAX_STACK_SIZE = 256
var MAX_LOCALS_SIZE = 256

type INTER_RESULT int

const (
	INTER_ILLEGAL INTER_RESULT = iota

	INTER_OK
	INTER_COMPILE_ERROR
	INTER_RUNTIME_ERROR
)

var valueTypeMap = map[OpKey]value.VALUE_TYPE{
	OpKey{a: value.VT_FLOAT, b: value.VT_FLOAT}: value.VT_FLOAT,
	OpKey{a: value.VT_INT, b: value.VT_FLOAT}:   value.VT_FLOAT,
	OpKey{a: value.VT_FLOAT, b: value.VT_INT}:   value.VT_FLOAT,
	OpKey{a: value.VT_INT, b: value.VT_INT}:     value.VT_INT,
}

type VM struct {
	chunk        *chunk.Chunk
	ip           *interface{}
	counter      int
	vstack       stack.Stack
	valueTypeMap map[OpKey]value.VALUE_TYPE
	globals      LookupTable
	strings      LookupTable
	current      parser.Compiler
}

type LookupTable struct {
	_map map[string]value.Value
}

func (l *LookupTable) findObj(o string) *value.Value {
	if value, ok := l._map[o]; ok {
		return &value
	}
	return nil
}

func (l *LookupTable) deleteObj(o string) {
	delete(l._map, o)
}

type OpKey struct {
	a value.VALUE_TYPE
	b value.VALUE_TYPE
}

func (vm *VM) InitVM() {
	vm.globals = LookupTable{
		_map: make(map[string]value.Value),
	}
	vm.strings = LookupTable{
		_map: make(map[string]value.Value),
	}
	vm.vstack = stack.Stack{
		Sarray: make([]value.Value, MAX_STACK_SIZE),
		Top:    -1,
	}
	vm.valueTypeMap = valueTypeMap
}

func (vm *VM) ResetStack() {
	vm.vstack = stack.Stack{
		Sarray: make([]value.Value, MAX_STACK_SIZE),
	}
}

func (vm *VM) FreeVM() {
	vm.vstack = stack.Stack{}
	// free objects
}

func (vm *VM) Move() interface{} {
	if vm.counter >= len(vm.chunk.Code) {
		return nil
	}
	vm.ip = &vm.chunk.Code[vm.counter]
	vm.counter++
	return *vm.ip
}

func (vm *VM) ReadConstant() value.Value {
	vm.Move()
	index := (*vm.ip).(uint)
	return vm.chunk.Constants.Values[index]
}

func (vm *VM) binaryOP(op string) bool {
	b := vm.vstack.Pop()
	a := vm.vstack.Pop()

	if !value.IsSameType(a.VT, b.VT) {
		vt, found := vm.valueTypeMap[OpKey{a: a.VT, b: b.VT}]
		if !found {
			return false
		}
		a, b = value.ConvertToExpectedType2(a, b, vt)
	}

	switch op {
	case "+":
		vm.vstack.Push(value.Add(&a, &b))
	case "-":
		vm.vstack.Push(value.Sub(&a, &b))
	case "/":
		vm.vstack.Push(value.Divide(&a, &b))
	case "*":
		vm.vstack.Push(value.Multiply(&a, &b))
	case ">":
		vm.vstack.Push(value.Greater(&a, &b))
	case "<":
		vm.vstack.Push(value.Less(&a, &b))
	}
	return true
}

func (v *VM) StackTrace() {
	fmt.Println("== Stack Trace ==")
	fmt.Println("[")
	for i := 0; i < v.vstack.Top+1; i++ {
		fmt.Printf("%d ", i)
		value.PrintValue(v.vstack.Sarray[i])
	}
	fmt.Println("]")
	fmt.Printf("== End Stack Trace ==\n\n")
}

func (vm *VM) run() INTER_RESULT {
	for {
		instruct := vm.Move()
		switch instruct {
		case codes.INSTRUC_CONSTANT:
			constant := vm.ReadConstant()
			vm.vstack.Push(constant)
		case codes.INSTRUC_NIL:
			vm.vstack.Push(value.New("", value.VT_NIL))
		case codes.INSTRUC_TRUE:
			vm.vstack.Push(value.NewBool(true))
		case codes.INSTRUC_FALSE:
			vm.vstack.Push(value.NewBool(false))
		case codes.INSTRUC_ERR:
			return INTER_RUNTIME_ERROR
		case codes.INSTRUC_NOT:
			_v, err := vm.vstack.Peek(0)
			if !value.IsBooleanType(_v.VT) || err != nil {
				return INTER_RUNTIME_ERROR
			}
			vm.vstack.Push(value.Negate(vm.vstack.Pop()))
		case codes.INSTRUC_NEGATE:
			a, err := vm.vstack.Peek(0)
			if !value.IsNumberType(a.VT) || err != nil {
				return INTER_RUNTIME_ERROR
			}
			vm.vstack.Push(value.Negate(vm.vstack.Pop()))
		case codes.INSTRUC_EQUAL:
			b := vm.vstack.Pop()
			a := vm.vstack.Pop()
			if !value.IsSameType(a.VT, b.VT) {
				vt, found := vm.valueTypeMap[OpKey{a: a.VT, b: b.VT}]
				if !found {
					return INTER_RUNTIME_ERROR
				}
				a, b = value.ConvertToExpectedType2(a, b, vt)
			}
			vm.vstack.Push(value.Equal(&a, &b))
		case codes.INSTRUC_ADDITION:
			if !vm.binaryOP("+") {
				return INTER_RUNTIME_ERROR
			}
		case codes.INSTRUC_SUBSTRACT:
			if !vm.binaryOP("-") {
				return INTER_RUNTIME_ERROR
			}
		case codes.INSTRUC_MULTIPLY:
			if !vm.binaryOP("*") {
				return INTER_RUNTIME_ERROR
			}
		case codes.INSTRUC_DIVIDE:
			if !vm.binaryOP("/") {
				return INTER_RUNTIME_ERROR
			}
		case codes.INSTRUC_GREATER:
			a, _ := vm.vstack.Peek(0)
			if !value.IsNumberType(a.VT) {
				return INTER_RUNTIME_ERROR
			}
			vm.binaryOP(">")
		case codes.INSTRUC_LESS:
			a, _ := vm.vstack.Peek(0)
			if !value.IsNumberType(a.VT) {
				return INTER_RUNTIME_ERROR
			}
			vm.binaryOP("<")
		case codes.INSTRUC_DECL_GLOBAL:
			cnst := vm.ReadConstant()
			declName := value.AsString(&cnst)
			v, _ := vm.vstack.Peek(0)
			_, found := vm.globals._map[*declName]
			if found {
				fmt.Println("Variable already declared", *declName)
				return INTER_RUNTIME_ERROR
			}
			vm.globals._map[*declName] = v
		case codes.INSTRUC_SET_DECL_GLOBAL:
			cnst := vm.ReadConstant()
			declName := value.AsString(&cnst)
			v, _ := vm.vstack.Peek(0)
			vm.globals._map[*declName] = v
		case codes.INSTRUC_GET_DECL_GLOBAL:
			cnst := vm.ReadConstant()
			declName := value.AsString(&cnst)
			v, found := vm.globals._map[*declName]
			if !found {
				fmt.Println("Variable not declared", *declName)
				return INTER_RUNTIME_ERROR
			}
			vm.vstack.Push(v)
		case codes.INSTRUC_PRINT:
			value.PrintValue(vm.vstack.Pop())
			fmt.Printf("\n")
		case codes.INSTRUC_POP:
			peek, err := vm.vstack.Peek(0)
			if err == nil {
				vm.vstack.Pop()
				value.PrintValue(peek)
			}
			fmt.Printf("\n")
		case codes.INSTRUC_RETURN:
			return INTER_OK
		}
		//vm.StackTrace()
	}
}

func Compile(source string, chk *chunk.Chunk) INTER_RESULT {
	lex := lexer.Init(source)
	comp := parser.Compiler{Locals: make([]parser.Local, MAX_LOCALS_SIZE)}
	p := parser.Init(lex, chk, &comp)

	p.Advance()
	for !p.Match(token.EOF) {
		p.Decl()
		// temp
		if p.Perror {
			break
		}
	}
	p.EndCompile()

	if p.Perror {
		return INTER_COMPILE_ERROR
	}
	return INTER_OK
}

func (vm *VM) Interpret(source string) INTER_RESULT {
	chk := chunk.Chunk{}

	if Compile(source, &chk) == INTER_COMPILE_ERROR {
		// parser.ppanic = true
		// parser.perror = true
		return INTER_COMPILE_ERROR
	}

	/* DEBUG */
	/*
		if vm.debug {
		}
	*/

	chunk.DissasChunk(&chk, "INSTRUCT")

	if len(chk.Code) != 0 {
		/* INIT START */
		vm.chunk = &chk
		vm.counter = 0
		vm.ip = &vm.chunk.Code[vm.counter]
		/* INIT END */
		return vm.run()
	}
	return INTER_OK
}
