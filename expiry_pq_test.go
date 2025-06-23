package main

import (
	"container/heap"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBasicExpirePQ(t *testing.T) {
	var h heap.Interface
	h = MakeExpirePQ(6)
	heap.Init(h)

	e := makeEntryAt(1, 1)
	heap.Push(h, e)

	e = makeEntryAt(2, 2)
	heap.Push(h, e)

	e = makeEntryAt(6, 3)
	heap.Push(h, e)

	e = makeEntryAt(3, 4)
	heap.Push(h, e)

	h.Push(makeEntryAt(4, 5))

	e = heap.Pop(h).(*entry)
	assert.Equal(t, 1, int(e.value))

	e = heap.Pop(h).(*entry)
	assert.Equal(t, 2, int(e.value))

	e = heap.Pop(h).(*entry)
	assert.Equal(t, 4, int(e.value))

	e = heap.Pop(h).(*entry)
	assert.Equal(t, 5, int(e.value))

	e = heap.Pop(h).(*entry)
	assert.Equal(t, 3, int(e.value))

	a := heap.Pop(h)
	assert.Nil(t, a)

}
