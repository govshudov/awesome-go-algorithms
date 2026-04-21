// Package faninout provides generic fan-out / fan-in helpers over channels.
package faninout

import (
	"context"
	"sync"
)

// FanOut distributes values from in across n output channels, round-robin.
// Each output channel is closed when in is closed or ctx is cancelled.
func FanOut[T any](ctx context.Context, in <-chan T, n int) []<-chan T {
	if n < 1 {
		n = 1
	}
	outs := make([]chan T, n)
	for i := range outs {
		outs[i] = make(chan T)
	}

	go func() {
		defer func() {
			for _, c := range outs {
				close(c)
			}
		}()
		idx := 0
		for {
			select {
			case <-ctx.Done():
				return
			case v, ok := <-in:
				if !ok {
					return
				}
				select {
				case <-ctx.Done():
					return
				case outs[idx] <- v:
					idx = (idx + 1) % n
				}
			}
		}
	}()

	result := make([]<-chan T, n)
	for i, c := range outs {
		result[i] = c
	}
	return result
}

// FanIn merges multiple input channels into a single output channel. The
// output channel is closed after all inputs are drained or ctx is cancelled.
func FanIn[T any](ctx context.Context, ins ...<-chan T) <-chan T {
	out := make(chan T)
	var wg sync.WaitGroup
	wg.Add(len(ins))
	for _, in := range ins {
		go func(in <-chan T) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case v, ok := <-in:
					if !ok {
						return
					}
					select {
					case <-ctx.Done():
						return
					case out <- v:
					}
				}
			}
		}(in)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
