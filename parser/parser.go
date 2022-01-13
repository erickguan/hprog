package parser

import (
	"github.com/badc0re/hprog/chunk"
	"github.com/badc0re/hprog/vm"
)

func TestParser(expression string) {
	/*
		lex := lexer.Init(expression)

		stck := stack.Stack{}
		for tkn := range lex.Consume() {
			fmt.Println(tkn)
			stck.Push(tkn)
		}
	*/
	v := vm.VM{}
	v.InitVM()
	chk := chunk.Chunk{}

	id := chk.AddVariable(123)
	chk.WriteChunk(chunk.OP_CONSTANT, 1)
	chk.WriteChunk(id, 1)

	id2 := chk.AddVariable(456)
	chk.WriteChunk(chunk.OP_CONSTANT, 1)
	chk.WriteChunk(id2, 1)

	chk.WriteChunk(chunk.OP_MULTIPLY, 1)
	chk.WriteChunk(chunk.OP_NEGATE, 1)
	chk.WriteChunk(chunk.OP_RETURN, 1)

	//dissasChunk(&chk, "simple instruction")

	v.Interpret(&chk)
	v.FreeVM()
	// freeChunk(&chk)
	// freeChunk(&chk)
}
