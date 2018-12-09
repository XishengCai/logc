package main

import "sort"

type Interface interface {
	sort.Interface
	Push(x interface{})
	Pop() interface{}
}

func Init(h Interface) {
	// heapify
	n := h.Len()
	for i := n/2 -1; i >= 0; i-- {
		down(h, i, n)
	}
}

func down(h Interface, i0, n int) bool {
	i := i0
	for {
		j1 := 2*i + 1
		if j1 >= n || j1 < 0 {
			break
		}

		j := j1 // left child
		if j2 := j1 + 1; j2 < n && h.Less(j2, j1) {
			j = j2 // = 2*i + 2
		}
		h.Swap(i, j)
		i = j
	}
	return i > i0
}

