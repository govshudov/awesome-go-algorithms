// Example program exercising a handful of packages in this repo.
// Run with: go run ./examples
package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/govshudov/awesome-go-algorithms/concurrency/circuitbreaker"
	"github.com/govshudov/awesome-go-algorithms/concurrency/workerpool"
	"github.com/govshudov/awesome-go-algorithms/datastructures/lrucache"
	"github.com/govshudov/awesome-go-algorithms/datastructures/trie"
)

func main() {
	demoLRU()
	demoTrie()
	demoWorkerPool()
	demoCircuitBreaker()
}

func demoLRU() {
	fmt.Println("== LRU cache ==")
	cache := lrucache.New[string, int](3)
	cache.Put("a", 1)
	cache.Put("b", 2)
	cache.Put("c", 3)
	cache.Get("a")    // "a" becomes most recent
	cache.Put("d", 4) // evicts "b"

	for _, key := range []string{"a", "b", "c", "d"} {
		if v, ok := cache.Get(key); ok {
			fmt.Printf("  %s -> %d\n", key, v)
		} else {
			fmt.Printf("  %s -> (evicted)\n", key)
		}
	}
}

func demoTrie() {
	fmt.Println("== Trie ==")
	t := trie.New()
	for _, w := range []string{"go", "golang", "gopher", "rust"} {
		t.Insert(w)
	}
	fmt.Printf("  Search(go)         = %v\n", t.Search("go"))
	fmt.Printf("  Search(goph)       = %v\n", t.Search("goph"))
	fmt.Printf("  HasPrefix(goph)    = %v\n", t.HasPrefix("goph"))
	fmt.Printf("  HasPrefix(python)  = %v\n", t.HasPrefix("python"))
}

func demoWorkerPool() {
	fmt.Println("== Worker pool ==")
	pool := workerpool.New(3)
	defer pool.Stop()

	var wg sync.WaitGroup
	for i := range 5 {
		wg.Add(1)
		i := i
		_ = pool.Submit(context.Background(), func(ctx context.Context) {
			defer wg.Done()
			time.Sleep(10 * time.Millisecond)
			fmt.Printf("  finished job %d\n", i)
		})
	}
	wg.Wait()
}

func demoCircuitBreaker() {
	fmt.Println("== Circuit breaker ==")
	cb := circuitbreaker.New(circuitbreaker.Config{
		FailureThreshold: 2,
		ResetTimeout:     50 * time.Millisecond,
	})

	flaky := errors.New("flaky")
	for i := range 4 {
		err := cb.Do(func() error {
			if i < 2 {
				return flaky
			}
			return nil
		})
		fmt.Printf("  call %d: state=%s err=%v\n", i, cb.State(), err)
	}
}
