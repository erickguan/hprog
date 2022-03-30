package stack

import "github.com/badc0re/hprog/value"

type (
	Node struct {
		value value.Value
		prev  *Node
	}
	Stack struct {
		top *Node
	}
)

func (stack *Stack) Push(value value.Value) {
	node := &Node{value, &Node{}}
	if stack.top == nil {
		stack.top = node
	} else {
		node.prev = stack.top
		stack.top = node
	}
}

func (stack *Stack) Pop() value.Value {
	value := stack.top.value
	prev := stack.top.prev
	stack.top = prev
	return value
}

func (stack *Stack) Peek() value.Value {
	value := stack.top.value
	prev := stack.top.prev
	stack.top = prev
	return value
}
