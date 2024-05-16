package main

import (
	"sync"
	"testing"
	"time"

	. "github.com/kofalt/go-memoize"
	"github.com/stretchr/testify/assert"
)

func TestConcurrency(t *testing.T) {
	counters := map[string]int{}
	var l sync.Mutex
	inc := func(k string) (int, error) {
		l.Lock()
		defer l.Unlock()
		counters[k]++

		return counters[k], nil
	}

	cache := NewMemoizer(90*time.Second, 10*time.Minute)

	concurrency := 10000
	var wg sync.WaitGroup
	wg.Add(concurrency)

	routine := func() {

		for _, k := range []string{"key1", "key2", "key3"} {
			k := k
			result, _, _ := Call(cache, k, func() (int, error) {
				return inc(k)
			})
			assert.Equal(t, 1, result)
		}

		wg.Done()
	}

	for range concurrency {
		go routine()
	}

	wg.Wait()

	expected := map[string]int{
		"key1": 1,
		"key2": 1,
		"key3": 1,
	}

	assert.Equal(t, expected, counters)
}
