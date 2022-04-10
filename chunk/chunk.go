package chunk

import (
	"fmt"

	"github.com/badc0re/hprog/codes"
	"github.com/badc0re/hprog/value"
)

type VarArray struct {
	Count  uint
	Values []value.Value
}

type Chunk struct {
	Count     uint
	Lines     []int
	Code      []interface{}
	Constants VarArray
}

func (c *Chunk) WriteChunk(code interface{}, line int) {
	c.Code = append(c.Code, code)
	c.Lines = append(c.Lines, line)
	c.Count++
}

func FreeChunk(chunk *Chunk) {

}

func OpInstruction(name string, offset uint) uint {
	fmt.Printf("%s\n", name)
	return offset + 1
}

func PrintConstant(name string, chunk *Chunk, offset uint) uint {
	constant := chunk.Code[offset+1].(uint)
	fmt.Printf("%-16s %d '", name, constant)
	value.PrintValue(chunk.Constants.Values[constant])
	fmt.Printf("'\n")
	return offset + 2
}

func DissasInstruction(chunk *Chunk, offset uint) uint {
	fmt.Printf("%04d ", offset)
	fmt.Printf("%4d ", chunk.Lines[offset])

	inst := chunk.Code[offset]
	switch inst {
	case codes.INSTRUC_CONSTANT:
		return PrintConstant("CONSTANT", chunk, offset)
	case codes.INSTRUC_ADDITION:
		return OpInstruction("INSTRUC_ADDITION", offset)
	case codes.INSTRUC_SUBSTRACT:
		return OpInstruction("INSTRUC_SUBSTRACT", offset)
	case codes.INSTRUC_MULTIPLY:
		return OpInstruction("INSTRUC_DIVIDE", offset)
	case codes.INSTRUC_DIVIDE:
		return OpInstruction("INSTRUC_DIVIDE", offset)
	case codes.INSTRUC_NEGATE:
		return OpInstruction("INSTRUC_NEGATE", offset)
	case codes.INSTRUC_NOT:
		return OpInstruction("INSTRUC_NOT", offset)
	case codes.INSTRUC_EQUAL:
		return OpInstruction("INSTRUC_EQUAL", offset)
	case codes.INSTRUC_GREATER:
		return OpInstruction("INSTRUC_GREATER", offset)
	case codes.INSTRUC_LESS:
		return OpInstruction("INSTRUC_LESS", offset)
	case codes.INSTRUC_FALSE:
		return OpInstruction("INSTRUC_FALSE", offset)
	case codes.INSTRUC_TRUE:
		return OpInstruction("INSTRUC_TRUE", offset)
	case codes.INSTRUC_NIL:
		return OpInstruction("INSTRUC_NIL", offset)
	case codes.INSTRUC_RETURN:
		return OpInstruction("INSTRUC_RETURN", offset)
	case codes.INSTRUC_PRINT:
		return OpInstruction("INSTRUC_PRINT", offset)
	case codes.INSTRUC_POP:
		return OpInstruction("INSTRUC_POP", offset)
	}
	// NOTE: should never reach!
	return 0
}

func DissasChunk(chunk *Chunk, name string) {
	fmt.Printf("\n")
	fmt.Printf("== %s == \n", name)

	offset := uint(0)
	//for pos, opr := range chunk.Code {
	for offset < chunk.Count {
		offset = DissasInstruction(chunk, offset)
		if offset == 0 {
			break
		}
	}
	fmt.Println()
}

func (c *Chunk) AddVariable(constant value.Value) uint {
	c.Constants.Values = append(c.Constants.Values, constant)
	c.Constants.Count++
	//constant.index = uint(chunk.Constants.count)
	return c.Constants.Count - 1
}
