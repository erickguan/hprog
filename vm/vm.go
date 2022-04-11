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

type INTER_RESULT int

const (
	INTER_ILLEGAL INTER_RESULT = iota

	INTER_OK
	INTER_COMPILE_ERROR
	INTER_RUNTIME_ERROR
)

type VM struct {
	chunk        *chunk.Chunk
	ip           *interface{}
	counter      int
	vstack       stack.Stack
	valueTypeMap map[OpKey]value.VALUE_TYPE
	objects      *value.ObjCtr
}

type OpKey struct {
	a value.VALUE_TYPE
	b value.VALUE_TYPE
}

func (vm *VM) InitVM() {
	vm.objects = nil
	vm.vstack = stack.Stack{
		Sarray: make([]value.Value, MAX_STACK_SIZE),
		Top:    -1,
	}
	valueTypeMap := map[OpKey]value.VALUE_TYPE{
		OpKey{a: value.VT_FLOAT, b: value.VT_FLOAT}: value.VT_FLOAT,
		OpKey{a: value.VT_INT, b: value.VT_FLOAT}:   value.VT_FLOAT,
		OpKey{a: value.VT_FLOAT, b: value.VT_INT}:   value.VT_FLOAT,
		OpKey{a: value.VT_INT, b: value.VT_INT}:     value.VT_INT,
		// OpKey{a: value.VT_STRING, b: value.VT_STRING}:     value.VT_STRING,
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

func (vm *VM) binaryOP(op string) INTER_RESULT {
	b := vm.vstack.Pop()
	a := vm.vstack.Pop()

	if !value.IsSameType(a.VT, b.VT) {
		vt := vm.valueTypeMap[OpKey{a: a.VT, b: b.VT}]
		a, b = value.ConvertToExpectedType2(a, b, vt)
	}
	/*
		// allow "+" for strings
		if !(value.IsNumberType(a.VT) &&
			value.IsNumberType(b.VT)) {
			// works!
			return INTER_RUNTIME_ERROR
		}
	*/

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
	return INTER_OK
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
		case codes.INSTRUC_NOT:
			_v, err := vm.vstack.Peek(0)
			if !value.IsBooleanType(_v.VT) || err != nil {
				// error
				return INTER_RUNTIME_ERROR
			}
			vm.vstack.Push(value.Negate(vm.vstack.Pop()))
		case codes.INSTRUC_NEGATE:
			a, err := vm.vstack.Peek(0)
			if !value.IsNumberType(a.VT) || err != nil {
				// error
				return INTER_RUNTIME_ERROR
			}
			vm.vstack.Push(value.Negate(vm.vstack.Pop()))
		case codes.INSTRUC_EQUAL:
			b := vm.vstack.Pop()
			a := vm.vstack.Pop()
			if !value.IsSameType(a.VT, b.VT) {
				vt := vm.valueTypeMap[OpKey{a: a.VT, b: b.VT}]
				a, b = value.ConvertToExpectedType2(a, b, vt)
			}
			vm.vstack.Push(value.Equal(&a, &b))
		case codes.INSTRUC_ADDITION:
			vm.binaryOP("+")
		case codes.INSTRUC_SUBSTRACT:
			vm.binaryOP("-")
		case codes.INSTRUC_MULTIPLY:
			vm.binaryOP("*")
		case codes.INSTRUC_DIVIDE:
			vm.binaryOP("/")
		case codes.INSTRUC_GREATER:
			a, _ := vm.vstack.Peek(0)
			if !value.IsNumberType(a.VT) {
				// error
				return INTER_RUNTIME_ERROR
			}
			vm.binaryOP(">")
		case codes.INSTRUC_LESS:
			a, _ := vm.vstack.Peek(0)
			if !value.IsNumberType(a.VT) {
				// error
				return INTER_RUNTIME_ERROR
			}
			vm.binaryOP("<")
		case codes.INSTRUC_PRINT:
			value.PrintValue(vm.vstack.Pop())
			fmt.Printf("\n")
			// NOT NEEDED!
			//return INTER_OK
		/*
			case codes.INSTRUC_POP:
				pop := vm.vstack.Pop()
				value.PrintValue(pop)
				fmt.Printf("\n")
		*/
		case codes.INSTRUC_RETURN:
			return INTER_OK
		}
		vm.StackTrace()
	}
}

func Compile(source string, chk *chunk.Chunk) INTER_RESULT {
	lex := lexer.Init(source)
	p := parser.Init(lex, chk)

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
