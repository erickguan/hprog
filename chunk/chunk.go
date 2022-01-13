package chunk

import (
	"fmt"

	"github.com/badc0re/hprog/value"
)

type OP_CODE int

const (
	OP_ILLEGAL OP_CODE = iota

	OP_ADD
	OP_SUBSTRACT
	OP_MULTIPLY
	OP_DIVIDE

	OP_CONSTANT
	OP_NEGATE

	OP_RETURN
)

type VarArray struct {
	Count  uint
	Values [1024]value.Value
}

type Chunk struct {
	Count     uint
	Lines     []uint
	Code      []interface{}
	Constants VarArray
}

func (c *Chunk) WriteChunk(code interface{}, line uint) {
	c.Code = append(c.Code, code)
	c.Lines = append(c.Lines, line)
	c.Count++
}

func FreeChunk(chunk *Chunk) {

}

func PrintConstant(name string, chunk *Chunk, offset uint) uint {
	constant := chunk.Code[offset+1].(uint)
	fmt.Printf("%-16s %d '", "OP_CONSTANT", constant)
	fmt.Printf("%g", chunk.Constants.Values[constant])
	fmt.Printf("'\n")
	return offset + 2
}

func DissasInstruction(chunk *Chunk, offset uint) uint {
	fmt.Printf("%04d ", offset)
	if offset > 0 && chunk.Lines[offset] == chunk.Lines[offset-1] {
		fmt.Printf("    | ")
	} else {
		fmt.Printf("%4d ", chunk.Lines[offset])
	}

	inst := chunk.Code[offset]
	switch inst {
	case OP_CONSTANT:
		return PrintConstant("OP_CONSTANT", chunk, offset)
	case OP_RETURN:
		fmt.Println("OP_RETURN")
		return offset + 1
	}
	// NOTE: should never reach!
	return 0
}

func DissasChunk(chunk *Chunk, name string) {
	fmt.Printf("== %s == \n", name)

	offset := uint(0)
	//for pos, opr := range chunk.Code {
	for offset < chunk.Count {
		offset = DissasInstruction(chunk, offset)
	}
}

func (c *Chunk) AddVariable(constant value.Value) uint {
	c.Constants.Values[c.Constants.Count] = constant
	c.Constants.Count++
	//constant.index = uint(chunk.Constants.count)
	return c.Constants.Count - 1
}
