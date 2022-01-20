package parser

import (
	"github.com/badc0re/hprog/vm"
)

func TestParser(expression string) {

	v := vm.VM{}
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

	v.Interpret("1+2+3")
	v.FreeVM()
	// freeChunk(&chk)
	// freeChunk(&chk)
}
