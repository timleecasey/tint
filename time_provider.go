package main

import "time"

type TimeProvider interface {
	NowInMillis() int64
}

type WallTime struct {
}

func (wt *WallTime) NowInMillis() int64 {
	return time.Now().UnixMilli()
}

type IncTimer struct {
	tick int64
}

//
// IncTimer returns a time incremented every time.
//
func (it *IncTimer) NowInMillis() int64 {
	it.tick++
	return it.tick
}


