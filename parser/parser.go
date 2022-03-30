package parser

import (
	"fmt"

	"github.com/badc0re/hprog/vm"
)

func TestParser(expression string) {

	v := vm.VM{}
	v.InitVM()
	/*
		v.InitVM()
		chk := chunk.Chunk{}

		id := chk.AddVariable(123)
		chk.WriteChunk(codes.INSTRUC_CONSTANT, 1)
		chk.WriteChunk(id, 1)

		id2 := chk.AddVariable(456)
		chk.WriteChunk(codes.INSTRUC_CONSTANT, 1)
		chk.WriteChunk(id2, 1)

		chk.WriteChunk(codes.INSTRUC_ADDITION, 1)

		//chk.WriteChunk(codes.OP_NEGATE, 1)

		id3 := chk.AddVariable(1000)
		chk.WriteChunk(codes.INSTRUC_CONSTANT, 1)
		chk.WriteChunk(id3, 1)

		chk.WriteChunk(codes.INSTRUC_ADDITION, 1)

		chk.WriteChunk(codes.INSTRUC_RETURN, 1)

		chunk.DissasChunk(&chk, "simple instruction")
	*/

	status := v.Interpret(expression)
	if status == vm.INTER_RUNTIME_ERROR {
		fmt.Println("Runtime error.")
	}
	v.FreeVM()
	// freeChunk(&chk)
	// freeChunk(&chk)
}
