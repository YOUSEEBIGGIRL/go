// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package heap provides heap operations for any type that implements
// heap.Interface. A heap is a tree with the property that each node is the
// minimum-valued node in its subtree.
//
// The minimum element in the tree is the root, at index 0.
//
// A heap is a common way to implement a priority queue. To build a priority
// queue, implement the Heap interface with the (negative) priority as the
// ordering for the Less method, so Push adds items while Pop removes the
// highest-priority item from the queue. The Examples include such an
// implementation; the file example_pq_test.go has the complete source.
//
package heap

import "sort"

// The Interface type describes the requirements
// for a type using the routines in this package.
// Any type that implements it may be used as a
// min-heap with the following invariants (established after
// Init has been called or if the data is empty or sorted):
//
//	!h.Less(j, i) for 0 <= i < h.Len() and 2*i+1 <= j <= 2*i+2 and j < h.Len()
//
// Note that Push and Pop in this interface are for package heap's
// implementation to call. To add and remove things from the heap,
// use heap.Push and heap.Pop.
type Interface interface {
	sort.Interface

	// Push 添加 x 到末尾
	Push(x interface{}) // add x as element Len()

	// Pop 弹出末尾元素并返回
	Pop() interface{}   // remove and return element Len() - 1.
}

// Init establishes the heap invariants required by the other routines in this package.
// Init is idempotent with respect to the heap invariants
// and may be called whenever the heap invariants may have been invalidated.
// The complexity is O(n) where n = h.Len().
//
// Init 使 h 堆化
func Init(h Interface) {
	// heapify
	n := h.Len()
	for i := n/2 - 1; i >= 0; i-- {
		down(h, i, n)
	}
}

// Push pushes the element x onto the heap.
// The complexity is O(log n) where n = h.Len().
func Push(h Interface, x interface{}) {
	h.Push(x)
	// 从下到上堆化
	up(h, h.Len()-1)
}

// Pop removes and returns the minimum element (according to Less) from the heap.
// The complexity is O(log n) where n = h.Len().
// Pop is equivalent to Remove(h, 0).
func Pop(h Interface) interface{} {
	n := h.Len() - 1
	// 堆顶是最小（或最大）元素，将其交换到末尾，pop 会移除该元素
	h.Swap(0, n)
	// FIXME down 的作用不确定，以下为猜测
	// down 会将第二小（或大）的元素移动到堆顶，便于下次操作
	down(h, 0, n)
	return h.Pop()
}

// Remove removes and returns the element at index i from the heap.
// The complexity is O(log n) where n = h.Len().
func Remove(h Interface, i int) interface{} {
	n := h.Len() - 1
	if n != i {
		h.Swap(i, n)
		if !down(h, i, n) {
			up(h, i)
		}
	}
	return h.Pop()
}

// Fix re-establishes the heap ordering after the element at index i has changed its value.
// Changing the value of the element at index i and then calling Fix is equivalent to,
// but less expensive than, calling Remove(h, i) followed by a Push of the new value.
// The complexity is O(log n) where n = h.Len().
func Fix(h Interface, i int) {
	if !down(h, i, h.Len()) {
		up(h, i)
	}
}

// up 代表上浮，从下到上堆化
// j 代表需要堆化的元素 index，up 会从该元素开始，不断向上调整堆
// 在 Push 中会调用 up
func up(h Interface, j int) {
	for {
		// j 的父节点
		i := (j - 1) / 2 // parent

		// i == j ：此时 i 是最后一个元素
		//
		// 在解释 !h.Less(j, i) 之前，先来看看 Less(i, j int) 接口的定义：
		//
		// 		func (m myheap) Less(i, j int) bool {
		//			return m[i] < m[j]	// 这代表实现的堆为最小堆
		// 		}
		//
		// 		func (m myheap) Less(i, j int) bool {
		//			return m[i] > m[j]	// 这代表实现的堆为最大堆
		// 		}
		//
		//
		// 1. 最小堆，如果 j（子节点）大于 i（父节点），情况如下图：
		//
		//						3	i
		//						  \
		//						   5  j
		//
		// 此时便不需要交换 i 和 j，!h.Less(j, i) 对应 !m[j] < m[i]，即 m[j] > m[i]
		//
		// 最大堆同理
		if i == j || !h.Less(j, i) {
			break
		}
		// 否则需要进行交换，还是以最小堆为例：
		//
		//						5	i
		//						  \
		//						   3  j
		//
		// 此时不满足 !h.Less(j, i)，此时 i（父节点）大于 j（子节点），
		// 不满足最小堆的定义，所以需要进行交换
		h.Swap(i, j)

		// up 代表上浮，所以更新 j 为其父节点 i
		j = i
	}
}

// down 代表下沉，从上到下堆化
// n := h.Len()，堆的长度
// i0 := n/2 - 1，当前 index，down 会从该元素开始向下堆化
// 对堆的 [i0, n] 这部分进行 down 操作
// 返回值 bool 代表 down 操作是否执行成功
// down 会忽略最后一个元素，比如以下情况：
//
// 			最小堆，执行 down 操作
//
//                 2
//            3         5	<- i 		此时父节点 i 在这里
//        99   10    233  1	 			j 初始为 j1，j1 = 2*i + 1 = [5] = 233
// 										按照最小堆的逻辑，应该将 j 更新为
// 										更小的右子节点j2（值为 1 ），然后交换 5 和 1
//
//			预料的情况
//
//                 2
//            3         1
//        99   10    233  5
//
//
//
//		     实际情况
//
// 代码中的逻辑：
//
// 		if j2 := j1 + 1; j2 < n && h.Less(j2, j1) {
//			j = j2
// 		}
//
// 此时 j2 := j1 + 1 = 6 = n，不满足 if 条件，所以不会更新 j 为 j2
// j 依然为 j1 233，j < j1，所以不执行 swap 操作，所以 down 的结果是：
//
//                 2
//            3         5
//        99   10    233  1
//
// 原因：Pop 中调用了 down，Pop 操作会移除最后一个元素，所以不能对
// 最后一个元素进行操作 这里就没有对最后元素 1 进行操作，所以执行 pop
// 操作时弹出元素为堆中最小元素 1，如果和预料的情况一样，对最后元素进行
// 操作，将导致 1 和 5 产生交换，最后元素为 5，这样 Pop 的结果就不正确了
func down(h Interface, i0, n int) bool {
	i := i0
	for {
		j1 := 2*i + 1          // j1 是 i 的左子节点
		// j1 > n 代表子节点不存在
		// j1 == n 代表子节点是堆中最后一个节点，down 操作会忽略最后一个元素
		if j1 >= n || j1 < 0 { // j1 < 0 after int overflow
			break
		}
		// j 用来保存子节点中较小（大）的那个，默认为左子节点
		j := j1 // left child

		// j2 := j1 + 1，这表示 j2 是 i 的右子节点
		// 这里是 j2 < n 而不是 j2 <= n，因为 down 操作会忽略最后一个元素
		// 如果是最小堆，则选出 j1（左子节点）和 j2（右子节点）中较小的那个
		// 如果是最大堆，则选出 j1（左子节点）和 j2（右子节点）中较大的那个
		if j2 := j1 + 1; j2 < n && h.Less(j2, j1) {
			// 满足条件，则更新 j 为右子节点
			j = j2 // = 2*i + 2  // right child
		}

		// 判断 j（子节点） 和 i（父节点）的关系
		// 如果是最小堆，!h.Less(j, i) 代表 j > i，即父节点小于较小的子节点，
		// 此时已经满足最小堆的特性了，直接 break
		// 如果是最大堆，!h.Less(j, i) 代表 j < i，与上面同理
		if !h.Less(j, i) {
			break
		}

		// 到这里说明 j（子节点） 和 i（父节点）不满足堆的特性，
		// 如果是最小堆，说明此时 i > j，需要交换
		// 如果是最大堆，说明此时 i < j，需要交换
		h.Swap(i, j)

		// down 代表下沉，更新 i 为其子节点 j，进行下一轮循环
		i = j
	}

	// 如果没有执行过 i = j，则不满足 i > i0 条件，
	// 这代表没有进行 down 操作
	return i > i0
}
