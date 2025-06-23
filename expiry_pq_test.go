package main

import (
	"container/heap"
	"testing"
)

func TestBasicExpirePQ(t *testing.T) {
	var h heap.Interface
	h = MakeExpirePQ(3)
	e := makeEntryAt(1, 1)
	heap.Init(h)

	h.Push(e)
	h.Push(makeEntryAt(2, 2))

	e = h.Pop().(*entry)
	e = h.Pop().(*entry)

}
