package main

import (
	"fmt"
)

const (
	//
	// LoadFactor
	// Is the factor used to cover overages between eviction by priority
	// and eviction by expiry.  This is set to cover unit testing.
	//
	LoadFactor = 2
)

type EntryPQ struct {
	entries  []*entry
	capacity int
	priority bool
	expiry   bool
}

func (pq *EntryPQ) String() string {
	return fmt.Sprintf("%v %v", pq.capacity, pq.entries)
}

func MakeExpirePQ(max int) *EntryPQ {
	return &EntryPQ{
		entries:  make([]*entry, max*LoadFactor),
		capacity: 0,
		expiry:   true,
	}
}

func MakePriorityPQ(max int) *EntryPQ {
	return &EntryPQ{
		entries:  make([]*entry, max*LoadFactor),
		capacity: 0,
		priority: true,
	}
}

func (pq *EntryPQ) Len() int {
	return pq.capacity
}

func (pq *EntryPQ) Swap(i, j int) {
	if i == -1 {
		return
	}
	if j == -1 {
		return
	}
	pq.entries[i], pq.entries[j] = pq.entries[j], pq.entries[i]
}

func (pq *EntryPQ) Push(a any) {
	pq.entries[pq.capacity] = a.(*entry)
	pq.capacity++
}

func (pq *EntryPQ) Pop() any {
	if pq.capacity <= 0 {
		return nil
	}
	pq.capacity--
	e := pq.entries[pq.capacity]
	pq.entries[pq.capacity] = nil
	return e
}

func (pq *EntryPQ) Less(i, j int) bool {
	if pq.expiry {
		// earliest at index 0
		less := pq.entries[i].expireTime < pq.entries[j].expireTime
		return less
	} else {
		// least priority at index 0
		less := int(pq.entries[i].priority) < int(pq.entries[j].priority)
		return less
	}
}

func (pq *EntryPQ) Peek() *entry {
	if pq.capacity == 0 {
		return nil
	}
	return pq.entries[0]
}
