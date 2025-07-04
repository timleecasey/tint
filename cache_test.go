package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestBasicExpire(t *testing.T) {

	c := NewCache(4, &WallTime{})

	c.Set("A", 1, 1, 1)
	c.Set("B", 2, 1, 2)
	c.Set("C", 3, 1, 3)
	c.Set("D", 4, 1, 4)

	var keys []string

	sleep(5)
	c.SetMaxItems(4)
	keys = c.Keys()
	assert.Equal(t, 4, len(keys))

	c.Set("E", 1, 1, 1)
	keys = c.Keys()
	fmt.Printf("%v\n", keys)
	assert.Equal(t, 1, len(keys))

	v, ok := c.Get("E")
	assert.True(t, ok)
	assert.Equal(t, 1, v)
}

func TestClear(t *testing.T) {
	c := NewCache(3, &WallTime{})
	c.Set("A", 1, 1, 100)
	c.Set("B", 1, 2, 100)
	c.Set("C", 1, 3, 100)
	var keys []string

	c.SetMaxItems(2)
	keys = c.Keys()
	assert.Equal(t, 2, len(keys))

	c.SetMaxItems(1)
	keys = c.Keys()
	assert.Equal(t, 1, len(keys))

	c.SetMaxItems(0)
	keys = c.Keys()
	assert.Equal(t, 0, len(keys))

}

func TestLRU(t *testing.T) {
	c := NewCache(3, &WallTime{})
	c.Set("A", 1, 1, 100)
	c.Set("B", 1, 1, 100)
	c.Set("C", 1, 1, 100)
	var keys []string

	c.Get("B")
	c.SetMaxItems(1)
	keys = c.Keys()
	assert.Equal(t, 1, len(keys))
	assert.Equal(t, "B", keys[0])

}

func TestShadow(t *testing.T) {
	log.Printf("START SHADOW\n")
	c := NewCache(10, &WallTime{})
	c.Set("A", 1, 10, 100)
	c.Set("A", 2, 20, 10)
	c.Set("A", 3, 99, 1)
	//c.Set("B", 3, 99, 50)
	//c.Set("C", 3, 99, 50)
	//c.Set("D", 3, 99, 50)
	//c.Set("E", 3, 99, 50)
	//c.Set("F", 3, 99, 50)
	var keys []string
	var ok bool
	var val int

	keys = c.Keys()
	assert.Equal(t, 1, len(keys))
	assert.Equal(t, "A", keys[0])

	val, ok = c.Get("A")
	assert.True(t, ok)
	assert.Equal(t, 3, val)

	sleep(2)
	log.Printf("AT 2s\n")
	val, ok = c.Get("A")
	assert.True(t, ok)
	assert.Equal(t, 2, val)

	sleep(10)
	log.Printf("AT 12s\n")

	val, ok = c.Get("A")
	assert.True(t, ok)
	assert.Equal(t, 1, val)

	log.Printf("DONE")
}

func TestEmptyExpire(t *testing.T) {
	c := NewCache(3, &WallTime{})
	c.Set("A", 1, 1, 1)
	c.Set("B", 2, 2, 1)
	c.Set("C", 3, 3, 3)
	var keys []string

	keys = c.Keys()
	assert.Equal(t, 3, len(keys))
	c.SetMaxItems(2)

	sleep(2)
	keys = c.Keys()
	assert.Equal(t, 2, len(keys))

	c.SetMaxItems(1)
	keys = c.Keys()
	assert.Equal(t, 1, len(keys))
	assert.Equal(t, "C", keys[0])

}

func TestGetExpired(t *testing.T) {
	c := NewCache(3, &WallTime{})
	c.Set("A", 1, 1, 1)
	c.Set("B", 1, 1, 1)
	c.Set("C", 1, 1, 10)
	var keys []string
	var ok bool
	var val int

	c.Get("B")
	keys = c.Keys()
	assert.Equal(t, 3, len(keys))

	sleep(2)
	val, ok = c.Get("B")
	assert.False(t, ok)
	assert.Equal(t, int(NO_KEY), val)
}

func TestPriorityOverflow(t *testing.T) {
	c := NewCache(3, &WallTime{})

	c.Set("A", 1, 1, 100)
	c.Set("B", 1, 2, 100)
	c.Set("C", 1, 3, 100)
	c.Set("D", 1, 4, 100)
	c.Set("E", 1, 5, 100)
	c.Set("F", 1, 6, 100)

	keys := c.Keys()
	assert.Equal(t, 3, len(keys))
	assert.Equal(t, "D", keys[0])
}


