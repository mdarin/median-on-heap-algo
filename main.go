package main

import (
	"container/heap"
	"errors"
	"fmt"
)

// An IntHeap is a min-heap of ints.
type IntHeap []int

func (h IntHeap) Len() int           { return len(h) }
func (h IntHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h IntHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *IntHeap) Push(x any) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*h = append(*h, x.(int))
}

func (h *IntHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]

	return x
}

func (h *IntHeap) Peek() any { return (*h)[0] }

type MaxHeap struct {
	IntHeap
}

func (h MaxHeap) Less(i, j int) bool { return h.IntHeap[i] > h.IntHeap[j] }

type MinHeap struct {
	IntHeap
}

func (h MinHeap) Less(i, j int) bool { return h.IntHeap[i] < h.IntHeap[j] }

type MedianHeap struct {
	Min *MinHeap
	Max *MaxHeap
}

func New() *MedianHeap {
	return &MedianHeap{
		Min: new(MinHeap),
		Max: new(MaxHeap),
	}
}

func (h *MedianHeap) Insert(value any) {
	v, ok := value.(int)
	if !ok {
		return
	}

	if h.Max.Len() < 2 {
		heap.Push(h.Max, value)

		return
	}

	top, ok := h.Max.Peek().(int)
	if !ok {
		return
	}

	if top >= v {
		heap.Push(h.Max, value)
	} else {
		heap.Push(h.Min, value)
	}

	// seze balancing

	if h.Max.Len() > h.Min.Len()+1 {
		vv := heap.Pop(h.Max)
		heap.Push(h.Min, vv)
	}

	if h.Min.Len() > h.Max.Len()+1 {
		vv := heap.Pop(h.Min)
		heap.Push(h.Max, vv)
	}
}

func (h *MedianHeap) Median() (int, error) {
	if h.Max.Len() == 0 && h.Min.Len() == 0 {
		return 0, errors.New("mean heap is empty")
	}

	if h.Max.Len() == h.Min.Len() {
		v1 := h.Max.Peek()
		v2 := h.Min.Peek()

		return (v1.(int) + v2.(int)) / 2, nil
	}

	if h.Max.Len() > h.Min.Len() {
		v := h.Max.Peek()

		return v.(int), nil
	} else {
		v := h.Min.Peek()

		return v.(int), nil
	}
}

// This example inserts several ints into an IntHeap, checks the minimum,
// and removes them in order of priority.
func main() {
	h := new(IntHeap)

	for _, v := range []int{2, 1, 5, 3, 7, -1, 6, 8, 11, 16} {
		heap.Push(h, v)
	}

	fmt.Printf("minimum: %d\n", h.Peek())

	mh := New()

	for h.Len() > 0 {
		v := heap.Pop(h)
		fmt.Printf("%d ", v.(int))
	}

	fmt.Println()

	for _, v := range []int{2, 1, 5, 3, 7, -1, 6, 8, 11, 16} {
		mh.Insert(v)
	}

	median, err := mh.Median()
	if err != nil {
		fmt.Println("err", err)

		return
	}

	fmt.Println("median", median)
}
