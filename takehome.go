package main

import (
	"cmp"
	"container/heap"
	"fmt"
	"log"
	"slices"
	"strings"
	"time"
)

/*
You can use any language.

Implement a cache which includes the following features:
  Expire Time - after which an entry in the cache is invalid
  Priority - lower priority entries should be evicted before higher priority entries
  LRU - within a given priority, the least recently used item should be evicted first

The cache should support these operations:
  Get: Get the value of the key if the key exists in the cache and is not expired.
  Set: Update or insert the value of the key with a priority value and expire time.
       This should never ever allow more items than maxItems to be in the cache.

The cache eviction strategy should be as follows:
  1. Evict an expired entry first.
  2. If there are no expired items to evict, evict the lowest priority entry.
  3. If there are multiple items with the same priority, evict the least
     recently used among them.

This data structure is expected to hold large datasets (>>10^6 items) and should be
as efficient as possible.
*/

type Priority uint8
type Value int

const (
	NO_KEY            = Value(-1)
	INITIAL_PRIORITY  = Priority(101)
	GREATEST_PRIORITY = Priority(100)
	LEAST_PRIORITY    = Priority(0)
)

var debug = false

type entry struct {
	expireTime int64    // expireTime is the time when the item after which an entry is considered expired
	op         uint64   // op is the nth operation used for LRU
	key        string   // key is used to go from expiry list to map eviction
	priority   Priority // priority of the entry
	value      Value    // the value of the entry
	created    int64    // purely for debugging, mark the created time so the millis become easier to read
}

func LowToHighExpire(a, b *entry) int {
	if a == nil && b == nil {
		return 0
	}
	if a == nil {
		return 1
	}
	if b == nil {
		return -1
	}
	return cmp.Compare(a.expireTime, b.expireTime)
}

func HiToLowPriority(a, b *entry) int {
	if a == nil && b == nil {
		return 0
	}
	if a == nil {
		return 1
	}
	if b == nil {
		return -1
	}
	return cmp.Compare(b.priority, a.priority)
}

func (e *entry) String() string {
	return fmt.Sprintf("(%v:%v:%v @ E:%v P:%v)", e.key, e.op, e.value, (e.expireTime-e.created)/1000, e.priority)
}

type PriorityExpiryCache struct {
	maxItems   int                 // maxItems is the maximum entries allowed in the cache
	itemCount  int                 // itemCount is the number of current entries in the cache
	curOp      uint64              // curOp will identify the order of entries and usage.
	expiryPq   heap.Interface      // expiryPq is a heap interface keeping the next expiry backed by an array
	priorities map[string][]*entry // priorities is a map of entires which are kept by priority
	timer      TimeProvider        // timer allows for a differentiable time provider, either wall clock or something less time consuming
}

func NewCache(capacity int, timer TimeProvider) *PriorityExpiryCache {
	var h heap.Interface
	h = MakeExpirePQ(capacity)
	heap.Init(h)
	return &PriorityExpiryCache{
		maxItems:   capacity,
		curOp:      0,
		expiryPq:   h,
		priorities: make(map[string][]*entry, 0),
		timer:      timer,
	}
}

func (c *PriorityExpiryCache) Keys() []string {
	keys := make([]string, len(c.priorities))

	index := 0
	for k, _ := range c.priorities {
		keys[index] = k
		index++
	}

	slices.Sort(keys)
	return keys
}

// Expected to be O(1)
func (c *PriorityExpiryCache) Get(key string) (int, bool) {
	key = strings.Trim(key, " \n\r\t")
	ofkey, ok := c.priorities[key]
	if !ok {
		return int(NO_KEY), false
	}

	now := c.timer.NowInMillis()

	//
	// If expired, expire items and refetch
	//
	if ofkey[0].expireTime < now {
		c.evictExpired(now)
		ofkey, ok = c.priorities[key]
		if !ok {
			return int(NO_KEY), false
		}
	}

	//
	// Update LRU
	//
	c.curOp++
	ofkey[0].op = c.curOp

	return int(ofkey[0].value), true
}

// Expected to be log(p) + log(e)
func (c *PriorityExpiryCache) Set(key string, value int, priority int, expireInSec int) {
	now := c.timer.NowInMillis()
	expireMs := int64(expireInSec*1000) + now
	c.curOp++
	i := &entry{
		expireTime: expireMs,
		priority:   Priority(priority),
		value:      Value(value),
		op:         c.curOp,
		key:        strings.Trim(key, " \n\r\t"),
		created:    time.Now().UnixMilli(), // fights with timer
	}

	//
	// Allow for duplication in the cache.  When there
	// is an eviction notice this should resolve then.
	//
	c.itemCount++

	if c.itemCount >= c.maxItems {
		c.evictItems()
	}

	heap.Push(c.expiryPq, i) // have to push from a heap

	ofKey, ok := c.priorities[key]
	if !ok {
		ofKey = make([]*entry, 1)
		ofKey[0] = i
		c.priorities[key] = ofKey
	} else {
		//
		// This seems overkill, or a waste of churn.
		//
		c.priorities[key] = append(ofKey, i)
		slices.SortFunc(c.priorities[key], HiToLowPriority)
	}
}

func (c *PriorityExpiryCache) SetMaxItems(maxItems int) {
	c.maxItems = maxItems
	if c.maxItems < c.itemCount {
		c.evictItems()
	}
}

// evictItems
// may evict items from the cache to make room for new ones. This is done
// in a tiered way.  First, if there are expired items, they are removed.
// If there are two priorities which shadow each other, the lower priority
// is removed.  It is import to note a lower priority item may stay in the
// cache if at least one higher priority item exists which will expire
// before the lower priority one.  If there is a tie for expire and priority
// then op is used to pick the most recent item to keep.
func (c *PriorityExpiryCache) evictItems() {

	c.evictExpired(c.timer.NowInMillis())

	if c.itemCount < c.maxItems {
		return
	}

	cnt := c.itemCount - c.maxItems
	for i := 0; i < cnt; i++ {
		c.evictPriorities()
	}
}

// evictPriorities
// Go through the items by priority and find a item to evict.
// Since we add/remove items one by one, this primarily operates
// by finding a low priority item and removing it.
func (c *PriorityExpiryCache) evictPriorities() {

	var cur *entry
	cur = nil
	lowestPri := INITIAL_PRIORITY
	lowestOp := c.curOp

	for _, v := range c.priorities {

		if len(v) > 0 {
			//
			// Select the lowest priority entry from a given key
			//
			e := v[len(v)-1]
			if e.priority < lowestPri {
				cur = e
				lowestPri = e.priority
				lowestOp = e.op
				//
				// Take the LRU if the priorities are the same.
				//
			} else if e.priority == lowestPri && e.op < lowestOp {
				cur = e
				lowestOp = e.op
			}
		}
	}

	if cur != nil {
		c.deleteItemFromMap(cur)
	}
}

// evictExpired
// Go through and remove items which are past now
func (c *PriorityExpiryCache) evictExpired(now int64) {

	var pq *EntryPQ
	pq = c.expiryPq.(*EntryPQ)
	e := pq.Peek()
	for e != nil && e.expireTime <= now {

		if debug {
			log.Printf("DEL EXP: %v\n", e)
		}

		e = heap.Pop(c.expiryPq).(*entry) // make sure and remove the same thing.
		c.deleteItemFromMap(e)

		e = pq.Peek()
	}

}

// deleteItemFromMap
// Given an expired item, remove it from the map.
// You may not find the item since it was removed by priority.
// An item is equal if it has the same key/priority/expire.
// Since there are different priorities and LRU considerations,
// these are taken over item value.
func (c *PriorityExpiryCache) deleteItemFromMap(toRm *entry) {

	//fmt.Printf("DEL %v\n", toRm)
	if entryList, ok := c.priorities[toRm.key]; ok {
		for i := range entryList {
			e := entryList[i]
			if e.priority == toRm.priority && e.expireTime == toRm.expireTime {
				c.itemCount--

				if debug {
					log.Printf("DEL PRI %v\n", e)
				}

				if len(entryList) == 1 {
					delete(c.priorities, toRm.key)
				} else {
					c.priorities[toRm.key] = append(entryList[0:i], entryList[i+1:]...)
				}
				break
			}
		}
	}
}

// whiskey tango foxtrot
func sleep(secs int) {
	st := time.Now().UnixMilli()
	en := st + int64(secs*1000)
	for time.Now().UnixMilli() < en {
		time.Sleep(1)
	}
}

func main() {
	//Example:
	c := NewCache(5, &WallTime{})
	// c.Set([key name], value, priority, expiryTime)
	c.Set("A", 1, 5, 100)
	c.Set("B", 2, 15, 1)
	c.Set("C", 3, 5, 10)
	c.Set("D", 4, 1, 15)
	c.Set("E", 5, 5, 150)
	c.Get("C")

	// Current time = 0
	c.SetMaxItems(5)
	fmt.Printf("KEYS %v\n", c.Keys())
	// space for 5 keys, all 5 items are included

	//fmt.Printf("TM %v \n", time.Now().UnixMilli())
	sleep(2)
	//fmt.Printf("TM %v \n", time.Now().UnixMilli())

	// Current time = 2
	c.SetMaxItems(4)
	//c.Keys() = ["A", "C", "D", "E"]
	fmt.Printf("KEYS %v -- A, C, D, E\n", c.Keys())

	// "B" is removed because it is expired.  expiry 3 < 5

	c.SetMaxItems(3)
	//c.Keys() = ["A", "C", "E"]
	fmt.Printf("KEYS %v -- A, C, E\n", c.Keys())

	// "D" is removed because it the lowest priority
	// D's expire time is irrelevant.

	c.SetMaxItems(2)
	//c.Keys() = ["C", "E"]
	fmt.Printf("KEYS %v -- C, E\n", c.Keys())

	// "A" is removed because it is least recently used."
	// A's expire time is irrelevant.

	c.SetMaxItems(1)
	//c.Keys() = ["C"]
	fmt.Printf("KEYS %v -- C\n", c.Keys())

	// "E" is removed because C is more recently used (due to the Get("C") event).
}
