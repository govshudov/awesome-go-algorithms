// Package pipeline composes channel-based stages into a pipeline that
// propagates context cancellation. Each stage receives values from an input
// channel and emits to an output channel.
package pipeline

import "context"

// Stage transforms an input stream into an output stream.
type Stage[In, Out any] func(ctx context.Context, in <-chan In) <-chan Out

// Source produces a channel of values from a slice, respecting ctx.
func Source[T any](ctx context.Context, values []T) <-chan T {
	out := make(chan T)
	go func() {
		defer close(out)
		for _, v := range values {
			select {
			case <-ctx.Done():
				return
			case out <- v:
			}
		}
	}()
	return out
}

// Map applies fn to each element from in and emits the result on the
// returned channel. The output channel is closed when in is closed or ctx
// is cancelled.
func Map[In, Out any](ctx context.Context, in <-chan In, fn func(In) Out) <-chan Out {
	out := make(chan Out)
	go func() {
		defer close(out)
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
				case out <- fn(v):
				}
			}
		}
	}()
	return out
}

// Filter emits values from in for which keep returns true.
func Filter[T any](ctx context.Context, in <-chan T, keep func(T) bool) <-chan T {
	out := make(chan T)
	go func() {
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				return
			case v, ok := <-in:
				if !ok {
					return
				}
				if !keep(v) {
					continue
				}
				select {
				case <-ctx.Done():
					return
				case out <- v:
				}
			}
		}
	}()
	return out
}

// Collect drains ch into a slice, stopping early if ctx is cancelled.
func Collect[T any](ctx context.Context, ch <-chan T) []T {
	var out []T
	for {
		select {
		case <-ctx.Done():
			return out
		case v, ok := <-ch:
			if !ok {
				return out
			}
			out = append(out, v)
		}
	}
}
