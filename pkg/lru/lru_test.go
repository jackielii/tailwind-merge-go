package lru

import (
	"fmt"
	"sync"
	"testing"
)

func TestEvictsLeastRecentlyUsedEntry(t *testing.T) {
	cache := Make(2)
	cache.Set("first", "1")
	cache.Set("second", "2")

	if got := cache.Get("first"); got != "1" {
		t.Fatalf("Get(%q) = %q; want %q", "first", got, "1")
	}

	cache.Set("third", "3")

	if got := cache.Get("second"); got != "" {
		t.Errorf("Get(%q) = %q; want eviction miss", "second", got)
	}
	if got := cache.Get("first"); got != "1" {
		t.Errorf("Get(%q) = %q; want %q", "first", got, "1")
	}
	if got := cache.Get("third"); got != "3" {
		t.Errorf("Get(%q) = %q; want %q", "third", got, "3")
	}
}

func TestConcurrentGetSet(t *testing.T) {
	const (
		goroutines = 32
		operations = 1_000
		keyCount   = 64
	)

	cache := Make(32)
	start := make(chan struct{})
	var wg sync.WaitGroup

	for worker := 0; worker < goroutines; worker++ {
		wg.Add(1)
		go func(worker int) {
			defer wg.Done()
			<-start

			for operation := 0; operation < operations; operation++ {
				key := fmt.Sprintf("key-%d", (worker+operation)%keyCount)
				value := "value-" + key
				cache.Set(key, value)
				if got := cache.Get(key); got != "" && got != value {
					t.Errorf("Get(%q) = %q; want %q or an eviction miss", key, got, value)
				}
			}
		}(worker)
	}

	close(start)
	wg.Wait()
}
