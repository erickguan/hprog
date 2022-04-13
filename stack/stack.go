package stack

import (
	"errors"

	"github.com/badc0re/hprog/value"
)

type Stack struct {
	Top      int
	capacity int
	Sarray   []value.Value
}

func (stack *Stack) Push(value value.Value) {
	stack.Top++
	stack.Sarray[stack.Top] = value
}

func (stack *Stack) Pop() (value.Value, error) {
	/*
		if stack.Top == -1 {
			return value.Value{}, errors.New("AA")
		}
	*/
	_r := stack.Sarray[stack.Top]
	stack.Top--
	return _r, nil
}

func (stack *Stack) Peek(distance int) (value.Value, error) {
	if stack.Top == -1 {
		return value.Value{}, errors.New("Stack is empty.")
	}
	return stack.Sarray[stack.Top-distance], nil
}

func (stack *Stack) Top1() int {
	return stack.Top
}
