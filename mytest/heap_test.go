package mytest

import (
	"container/heap"
	"fmt"
	"testing"
)

type myheap []int

func (m myheap) Len() int {
	return len(m)
}

func (m myheap) Less(i, j int) bool {
	return m[i] < m[j]
}

func (m myheap) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

func (m *myheap) Push(x interface{}) {
	*m = append(*m, x.(int))
}

func (m *myheap) Pop() interface{} {
	p := (*m)[m.Len()-1]
	*m = (*m)[:m.Len()-1]
	return p
}

func TestHeapPushPop(t *testing.T) {
	h := new(myheap)
	heap.Push(h, 5)
	heap.Push(h, 10)
	heap.Push(h, 2)
	heap.Push(h, 99)
	heap.Push(h, 3)
	heap.Push(h, 233)
	heap.Push(h, 1)
	heap.Pop(h)
}

func TestHeapInit(t *testing.T) {
	var h myheap
	h = append(h, 50, 2, 60, 3, 1)
	fmt.Println(h)

	heap.Init(&h)
	fmt.Println(h)
}



