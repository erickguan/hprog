package stack

type (
	Node struct {
		value interface{}
		prev  *Node
	}
	Stack struct {
		top *Node
	}
)

func (stack *Stack) Push(value interface{}) {
	node := &Node{value, &Node{}}
	if stack.top == nil {
		stack.top = node
	} else {
		node.prev = stack.top
		stack.top = node
	}
}

func (stack *Stack) Pop() interface{} {
	if stack.top == nil {
		return nil
	}
	value := stack.top.value
	prev := stack.top.prev
	stack.top = prev
	return value
}
