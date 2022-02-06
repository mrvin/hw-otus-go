package hw04lrucache

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"
)

type sTest struct {
	keyInput Key
	valInput int
	want     bool
}

type gTest struct {
	keyInput Key
	wantVal  interface{}
	want     bool
}

func testGet(t *testing.T, test gTest, c Cache) {
	val, ok := c.Get(test.keyInput)
	if val != test.wantVal || ok != test.want {
		t.Errorf("cache.Get(%q) = %v, %t; want: %v, %t", test.keyInput, val, ok, test.wantVal, test.want)
	}
}

func testSet(t *testing.T, test sTest, c Cache) {
	ok := c.Set(test.keyInput, test.valInput)
	if ok != test.want {
		t.Errorf("cache.Set(%q, %v) = %t; want: %t", test.keyInput, test.valInput, ok, test.want)
	}
}

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)
		tests := []gTest{
			{"aaa", nil, false},
			{"bbb", nil, false},
		}

		for _, test := range tests {
			testGet(t, test, c)
		}
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		testSet(t, sTest{"aaa", 100, false}, c)
		testSet(t, sTest{"bbb", 200, false}, c)

		testGet(t, gTest{"aaa", 100, true}, c)
		testGet(t, gTest{"bbb", 200, true}, c)

		testSet(t, sTest{"aaa", 300, true}, c)

		testGet(t, gTest{"aaa", 300, true}, c)
		testGet(t, gTest{"ccc", nil, false}, c)

	})

	t.Run("purge logic (capacity)", func(t *testing.T) {
		c := NewCache(3)
		tests := []sTest{
			{"aaa", 100, false},
			{"bbb", 200, false},
			{"ccc", 300, false},
			{"ddd", 200, false},
		}

		for _, test := range tests {
			testSet(t, test, c)
		}

		testGet(t, gTest{"aaa", nil, false}, c)
	})

	t.Run("purge logic (rarely used)", func(t *testing.T) {
		c := NewCache(3)

		testSet(t, sTest{"aaa", 100, false}, c)
		testSet(t, sTest{"bbb", 200, false}, c)
		testSet(t, sTest{"ccc", 300, false}, c)

		testGet(t, gTest{"aaa", 100, true}, c)
		testGet(t, gTest{"bbb", 200, true}, c)

		testSet(t, sTest{"aaa", 400, true}, c)
		testSet(t, sTest{"ddd", 500, false}, c)

		testGet(t, gTest{"ccc", nil, false}, c)
	})

	t.Run("Clear", func(t *testing.T) {
		c := NewCache(3)
		testsSet := []sTest{
			{"aaa", 100, false},
			{"bbb", 200, false},
			{"ccc", 300, false},
		}
		testsGet := []gTest{
			{"aaa", nil, false},
			{"bbb", nil, false},
			{"ccc", nil, false},
		}

		for _, test := range testsSet {
			testSet(t, test, c)
		}

		c.Clear()

		for _, test := range testsGet {
			testGet(t, test, c)
		}
	})
}

func TestCacheMultithreading(t *testing.T) {
	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}
