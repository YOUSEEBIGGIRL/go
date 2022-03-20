package mytest

import (
	"fmt"
	"testing"
)

type Stack[T any] struct {
	val []T
}

func NewStack[T any]() *Stack[T] {
	return &Stack[T]{}
}

func (s *Stack[T]) Push(v T) {
	s.val = append(s.val, v)
}

func (s *Stack[T]) Pop() T {
	if len(s.val) == 0 {
		panic("stack size is zero")
	}
	p := s.val[len(s.val)-1]
	s.val = s.val[:len(s.val)-1]
	return p
}

func (s *Stack[T]) Top() T {
	if len(s.val) == 0 {
		panic("stack size is zero")
	}
	return s.val[len(s.val)-1]
}

func printAny(val any) {
	fmt.Println(val)
}

func Test1(t *testing.T) {
	stack := NewStack[int]()

	stack.Push(1)
	stack.Push(2)
	fmt.Println(stack)

	stack.Pop()
	fmt.Println(stack.Top())

	stack.Pop()
	stack.Pop()
	fmt.Println(stack)

}

func Test2(t *testing.T) {
	printAny(1)
	printAny("1")
	printAny(struct {
		X, Y int64
	}{1, 2})
	printAny(1.1)
	printAny([]int{1, 2, 3})
}
