package tokeniser

// LRU stack
// A stack with a maximum size that discards the oldest
// item when the stack is full
type Stack[T any] struct {
	items []T
	top   uint
	size  uint
}

// NewStack creates a new stack with a maximum size
func NewStack[T any](size uint) *Stack[T] {
	return &Stack[T]{
		items: make([]T, size),
		size:  size,
	}
}

// Push adds an item to the stack discarding the oldest
// item if the stack is full
func (s *Stack[T]) Push(item T) {
	if s.top == s.size {
		for i := uint(1); i < s.size; i++ {
			s.items[i-1] = s.items[i]
		}
		s.top--
	}

	s.items[s.top] = item
	s.top++
}

// Pop removes the last item from the stack and returns it
// or panics if the stack is empty
func (s *Stack[T]) Pop() T {
	if s.top == 0 {
		panic("Stack is empty")
	}

	s.top--
	return s.items[s.top]
}

// Top returns the last item from the stack without removing it
func (s *Stack[T]) Top() T {
	if s.top == 0 {
		panic("Stack is empty")
	}

	return s.items[s.top-1]
}

// Length returns the number of items in the stack
func (s *Stack[T]) Length() uint {
	return s.top
}
