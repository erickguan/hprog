package parser

import (
	"fmt"
)

type OP_CODE int

const (
	OP_ILLEGAL OP_CODE = iota

	/*
		OPR_ADD
		OPR_DIVIDE
		OPR_MULTIPLY
		OPR_SUBSTRACT
	*/
	OP_CONSTANT
	OP_RETURN
)

type VT_CODE int

const (
	VT_ILLEGAL VT_CODE = iota

	/*
		OPR_ADD
		OPR_DIVIDE
		OPR_MULTIPLY
		OPR_SUBSTRACT
	*/
	VT_INT
	VT_RETURN
)

type Chunk struct {
	count     uint
	lines     []uint
	code      []interface{}
	constants VarArray
}

func writeChunk(chunk *Chunk, code interface{}, line uint) {
	chunk.code = append(chunk.code, code)
	chunk.lines = append(chunk.lines, line)
	chunk.count++
}

func freeChunk(chunk *Chunk) {

}

func printConstant(name string, chunk *Chunk, offset uint) uint {
	constant := chunk.code[offset+1].(uint)
	fmt.Printf("%-16s %d '", "OP_CONSTANT", constant)
	fmt.Printf("%g", chunk.constants.values[constant])
	fmt.Printf("'\n")
	return offset + 2
}

func printValue(value Value) {
	fmt.Printf("val: %g", value)
}

func dissasInstruction(chunk *Chunk, offset uint) uint {
	fmt.Printf("%04d ", offset)
	if offset > 0 && chunk.lines[offset] == chunk.lines[offset-1] {
		fmt.Printf("    | ")
	} else {
		fmt.Printf("%4d ", chunk.lines[offset])
	}

	inst := chunk.code[offset]
	switch inst {
	case OP_CONSTANT:
		return printConstant("OP_CONSTANT", chunk, offset)
	case OP_RETURN:
		fmt.Println("OP_RETURN")
		return offset + 1
	}
	// NOTE: should never reach!
	return 0
}

func dissasChunk(chunk *Chunk, name string) {
	fmt.Printf("== %s == \n", name)

	offset := uint(0)
	//for pos, opr := range chunk.code {
	for offset < chunk.count {
		offset = dissasInstruction(chunk, offset)
	}
}

/*
type Variable struct {
	vtype string
	value interface{}
}
*/

type VarArray struct {
	count  uint
	values [1024]Value
}

type Value float64

func addVariable(chunk *Chunk, constant Value) uint {
	chunk.constants.values[chunk.constants.count] = constant
	chunk.constants.count++
	//constant.index = uint(chunk.constants.count)
	return chunk.constants.count - 1
}

type VM struct {
	chunk *Chunk
	//ip      *interface{}
	counter int
}

type INTER_RESULT int

const (
	INTER_ILLEGAL INTER_RESULT = iota

	INTER_OK
	INTER_COMPILE_ERROR
	INTER_RUNTIME_ERROR
)

func initVM() {

}

func freeVM() {

}

func (vm *VM) move() int {
	//vm.ip = &vm.chunk.code[vm.counter]
	ctr := vm.counter
	vm.counter++
	return ctr
}

func readConstant(vm *VM) Value {
	index := vm.chunk.code[vm.move()].(uint)
	return vm.chunk.constants.values[index]
}

func run(vm *VM) INTER_RESULT {
	for {
		instruct := vm.chunk.code[vm.move()]
		switch instruct {
		case OP_CONSTANT:
			constant := readConstant(vm)
			printValue(constant)
			break
		case OP_RETURN:
			return INTER_OK
		}
	}
}

func interpret(vm *VM, chunk *Chunk) INTER_RESULT {
	vm.chunk = chunk
	vm.counter = 0
	//vm.ip = &vm.chunk.code[vm.counter]
	return run(vm)
}

func TestParser(expression string) {
	/*
		lex := lexer.Init(expression)

		stck := stack.Stack{}
		for tkn := range lex.Consume() {
			fmt.Println(tkn)
			stck.Push(tkn)
		}
	*/
	initVM()
	chunk := Chunk{}

	id := addVariable(&chunk, 123)
	writeChunk(&chunk, OP_CONSTANT, 1)
	writeChunk(&chunk, id, 1)
	writeChunk(&chunk, OP_RETURN, 1)

	//dissasChunk(&chunk, "simple instruction")

	vm := VM{}
	interpret(&vm, &chunk)
	freeVM()
	// freeChunk(&chunk)
	// freeChunk(&chunk)
}
