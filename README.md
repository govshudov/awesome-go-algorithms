# awesome-go-algorithms

A curated collection of generic data structures and concurrency patterns in idiomatic Go.

Zero external dependencies. Standard library only. Go 1.22+ generics throughout.

## Install

```bash
go get github.com/govshudov/awesome-go-algorithms/...
```

## Packages

### Data structures

| Package | Description |
| --- | --- |
| [`datastructures/linkedlist`](datastructures/linkedlist) | Generic doubly-linked list with O(1) push/pop at both ends |
| [`datastructures/stack`](datastructures/stack) | Generic slice-backed stack |
| [`datastructures/queue`](datastructures/queue) | Generic ring-buffer queue with amortized O(1) enqueue/dequeue |
| [`datastructures/heap`](datastructures/heap) | Generic binary heap parameterized by a comparator (use as min- or max-heap) |
| [`datastructures/trie`](datastructures/trie) | Rune-based trie supporting `Insert`, `Search`, `HasPrefix`, `Delete` |
| [`datastructures/lrucache`](datastructures/lrucache) | Generic O(1) LRU cache built on `container/list` |
| [`datastructures/unionfind`](datastructures/unionfind) | Disjoint-set with path compression and union by rank |

### Concurrency patterns

| Package | Description |
| --- | --- |
| [`concurrency/workerpool`](concurrency/workerpool) | Fixed-size, context-aware worker pool with graceful shutdown |
| [`concurrency/pipeline`](concurrency/pipeline) | Channel-based pipeline stages with cancellation propagation |
| [`concurrency/faninout`](concurrency/faninout) | Generic `FanOut` and `FanIn` helpers |
| [`concurrency/ratelimiter`](concurrency/ratelimiter) | Token-bucket rate limiter, goroutine-safe |
| [`concurrency/circuitbreaker`](concurrency/circuitbreaker) | Classic Closed/Open/HalfOpen circuit breaker state machine |

## Usage

### LRU cache

```go
import "github.com/govshudov/awesome-go-algorithms/datastructures/lrucache"

cache := lrucache.New[string, int](3)
cache.Put("a", 1)
cache.Put("b", 2)
if v, ok := cache.Get("a"); ok {
    fmt.Println(v) // 1
}
```

### Worker pool

```go
import "github.com/govshudov/awesome-go-algorithms/concurrency/workerpool"

pool := workerpool.New(4)
defer pool.Stop()

for i := 0; i < 10; i++ {
    i := i
    pool.Submit(context.Background(), func(ctx context.Context) {
        fmt.Println("job", i)
    })
}
```

### Circuit breaker

```go
import "github.com/govshudov/awesome-go-algorithms/concurrency/circuitbreaker"

cb := circuitbreaker.New(circuitbreaker.Config{
    FailureThreshold: 3,
    ResetTimeout:     5 * time.Second,
})

err := cb.Do(func() error {
    return callFlakyService()
})
```

## Running tests

```bash
go test ./...
go test -race ./concurrency/...
```

## Contributing

Pull requests welcome. Please keep packages dependency-free and include table-driven tests.

## License

MIT - see [LICENSE](LICENSE).
