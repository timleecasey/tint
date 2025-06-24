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

	e := makeExpEntryAt(1, 1)
	heap.Push(h, e)

	e = makeExpEntryAt(2, 2)
	heap.Push(h, e)

	e = makeExpEntryAt(6, 3)
	heap.Push(h, e)

	e = makeExpEntryAt(3, 4)
	heap.Push(h, e)

	h.Push(makeExpEntryAt(4, 5))

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
func TestBasicPriorityPQ(t *testing.T) {
	var h heap.Interface
	h = MakePriorityPQ(6)
	heap.Init(h)

	e := makePriEntry(1, 1)
	heap.Push(h, e)

	e = makePriEntry(2, 2)
	heap.Push(h, e)

	e = makePriEntry(6, 3)
	heap.Push(h, e)

	e = makePriEntry(3, 4)
	heap.Push(h, e)

	h.Push(makePriEntry(4, 5))

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
