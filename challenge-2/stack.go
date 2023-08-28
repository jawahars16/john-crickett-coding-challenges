package parser

import "container/list"

type Stack struct {
	list list.List
}

func NewStack() *Stack {
	return &Stack{
		list: *list.New(),
	}
}

func (s *Stack) Push(v any) {
	s.list.PushFront(v)
}

func (s *Stack) Pop() any {
	v := s.list.Front()
	if v != nil {
		s.list.Remove(v)
		return v.Value
	}
	return nil
}

func (s *Stack) Peek() any {
	e := s.list.Front()
	if e != nil {
		return e.Value
	}
	return nil
}
