package collections
import "container/list"

type Stack struct {
	values *list.List
}

func (s Stack) Size() int {
	return s.values.Len()
}


//create a new stack
func NewStack(size int) *Stack {
	if size < 0 {
		panic(IllegalArgumentError{"size must >= 0"})
	}
	stack := &Stack{}
	stack.values = list.New()
	return stack
}

// push a value into stack
func (s *Stack)Push(value interface{}) *Stack {
	s.values.PushBack(value)
	return s
}

// pop a value from stack,
// if has value to pop, bool will is true, otherwise is false
func (s *Stack)Pop() (interface{}, bool) {
	if s.values.Len() == 0 {
		return nil, false
	}

	e := s.values.Back()
	return s.values.Remove(e), true
}

// clear the stack
func (s *Stack)Clear() *Stack {
	s.values.Init()
	return s
}


// judge the stack if has elements
// return true if has no elements, otherwise return true
func (s *Stack)IsEmpty() bool {
	return s.values.Len() == 0
}

func (s *Stack) Elements() []interface{} {
	elements := make([]interface{}, 0, s.values.Len())

	for e := s.values.Front(); e != nil; e = e.Next() {
		elements = append(elements, e.Value)
	}

	return elements
}

