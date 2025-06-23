package main

type ExpirePQ struct {
	entries  []*entry
	capacity int
}

func MakeExpirePQ(max int) *ExpirePQ {
	return &ExpirePQ{
		entries:  make([]*entry, max),
		capacity: 0,
	}
}

func (pq *ExpirePQ) Len() int {
	return pq.capacity
}

func (pq *ExpirePQ) Swap(i, j int) {
	pq.entries[i], pq.entries[j] = pq.entries[j], pq.entries[j]
}

func (pq *ExpirePQ) Push(a any) {
	pq.entries[pq.capacity] = a.(*entry)
	pq.capacity++
}

func (pq *ExpirePQ) Pop() any {
	if pq.capacity <= 0 {
		return nil
	}
	pq.capacity--
	e := pq.entries[pq.capacity]
	pq.entries[pq.capacity] = nil
	return e
}

func (pq *ExpirePQ) Less(i, j int) bool {
	return pq.entries[i].expireTime < pq.entries[j].expireTime
}

func (pq *ExpirePQ) More(i, j int) bool {
	return pq.entries[i].expireTime > pq.entries[j].expireTime
}

func (pq *ExpirePQ) Peek() *entry {
	if pq.capacity == 0 {
		return nil
	}
	return pq.entries[pq.capacity-1]
}
