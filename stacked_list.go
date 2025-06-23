package main

import (
	"slices"
)

type StackedList interface {
	Append(*entry)
	Size() int
	Get(i int) *entry
	Delete(i int) *entry
}

type EntryChunk struct {
	entries   []*entry // entries are sorted by expire time
	highWater int      // highWater is the next index to fill with an entry
	lowWater  int      // lowWater is the first index of a valid entry
	sorted    bool
}

type ChunkedExpiryList struct {
	entries   []*EntryChunk // blocks of entries
	size      int           // the total number fo entries
	endBlock  int           // endBlock is the index of the block to append into
	blockSize int           // blockSize the size of blocks
}

func (sl *ChunkedExpiryList) MakeEntryBlock() *EntryChunk {
	eb := &EntryChunk{
		entries:   make([]*entry, sl.blockSize),
		highWater: 0,
		lowWater:  0,
		sorted:    false,
	}

	return eb
}

func MakeExpiryList(blockSize int) *ChunkedExpiryList {
	ret := &ChunkedExpiryList{
		entries:   make([]*EntryChunk, 1),
		endBlock:  0,
		size:      0,
		blockSize: blockSize,
	}
	ret.entries[0] = ret.MakeEntryBlock()

	return ret
}

func (eb *EntryChunk) appendToChunk(e *entry) {
	eb.entries[eb.highWater] = e
	eb.highWater++
	eb.sorted = false
}

func (eb *EntryChunk) handleSort() {
	if !eb.sorted {
		slices.SortFunc(eb.entries, LowToHighExpire)
		eb.sorted = true
	}
}

func (eb *EntryChunk) get(i int) *entry {
	eb.handleSort()
	return eb.entries[i+eb.lowWater]
}

func (eb *EntryChunk) delete() *entry {
	eb.handleSort()
	e := eb.entries[eb.lowWater]
	eb.entries[eb.lowWater] = nil
	eb.lowWater++
	return e
}

func (el *ChunkedExpiryList) Append(e *entry) {
	if el.entries[el.endBlock].highWater >= el.blockSize {
		el.entries = append(el.entries, el.MakeEntryBlock())
		el.endBlock++
	}
	eb := el.entries[el.endBlock]
	eb.appendToChunk(e)
	el.size++
}

func (el *ChunkedExpiryList) Size() int {
	return el.size
}

func (el *ChunkedExpiryList) Get(i int) *entry {
	blockIndex := el.indexToChunk(i)
	eb := el.entries[blockIndex]
	return eb.get(i % el.blockSize)
}

func (el *ChunkedExpiryList) Delete() *entry {
	eb := el.entries[0]
	toRm := eb.delete()
	if eb.lowWater >= eb.highWater {
		if len(el.entries) > 1 {
			el.entries = el.entries[1:]
			el.endBlock--
		} else {
			el.entries[0] = el.MakeEntryBlock()
		}
	}
	if toRm != nil {
		el.size--
	}
	return toRm
}

// indexToChunk
// Computes an index based on the count of each chunk
// This assumes the chunks are large enough to not dominate
// So for 1 million entries, a chunk size  of 10k would be 100 long.
// Maybe that is close enough to ~constant.
// If not, then assume the first chunk is partial and then assume the
// rest are full and compute an index value.
func (el *ChunkedExpiryList) indexToChunk(i int) int {
	count := i
	for i := range el.entries {
		eb := el.entries[i]
		if eb.highWater-eb.lowWater >= count {
			return i
		}
		count -= eb.highWater - eb.lowWater
	}
	return -1
}
