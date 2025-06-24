package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBasicStacked(t *testing.T) {

	sl := MakeExpiryList(3)
	e := makeExpEntryAt(1, 1)
	sl.Append(e)
	sl.Append(e)
	sl.Append(e)

	sl.Append(e)
	sl.Append(e)
	sl.Append(e)

	sl.Append(e)
	sl.Append(e)

	sz := sl.Size()
	assert.Equal(t, 8, sz)
	assert.Equal(t, 2, sl.endBlock)
	assert.Equal(t, 3, len(sl.entries))
	assert.Equal(t, 3, sl.blockSize)

	e = sl.Delete()
	e = sl.Delete()
	e = sl.Delete()

	e = sl.Delete()
	e = sl.Delete()
	e = sl.Delete()

	e = sl.Delete()
	e = sl.Delete()
	e = sl.Delete()
	assert.Nil(t, e)

	sz = sl.Size()
	assert.Equal(t, 0, sz)
	assert.Equal(t, 0, sl.endBlock)
	assert.Equal(t, 1, len(sl.entries))
	assert.Equal(t, 3, sl.blockSize)
}

func TestExpire(t *testing.T) {
	sl := MakeExpiryList(3)
	var e *entry

	e = makeExpEntryAt(2, 1)
	sl.Append(e)

	e = makeExpEntryAt(2, 2)
	sl.Append(e)

	e = makeExpEntryAt(2, 3)
	sl.Append(e)

	e = makeExpEntryAt(1, 4)
	sl.Append(e)

	e = sl.Delete()
	assert.Equal(t, int64(2), e.expireTime)

	e = sl.Delete()
	assert.Equal(t, int64(2), e.expireTime)

	e = sl.Delete()
	assert.Equal(t, int64(2), e.expireTime)

	e = sl.Delete()
	assert.Equal(t, int64(1), e.expireTime)

}

func makeExpEntryAt(expire int64, n int) *entry {
	e := &entry{
		expireTime: expire,
		op:         uint64(n),
		key:        fmt.Sprintf("k%v", n),
		priority:   Priority(n),
		value:      Value(n),
	}
	return e
}

func makePriEntry(priority int, n int) *entry {
	e := &entry{
		expireTime: 100,
		op:         uint64(n),
		key:        fmt.Sprintf("k%v", n),
		priority:   Priority(priority),
		value:      Value(n),
	}
	return e
}
