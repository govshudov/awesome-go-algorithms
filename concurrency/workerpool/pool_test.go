package workerpool

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

func TestPool_RunsTasks(t *testing.T) {
	p := New(4)
	defer p.Stop()

	var counter atomic.Int32
	const n = 100
	done := make(chan struct{}, n)

	for range n {
		err := p.Submit(context.Background(), func(ctx context.Context) {
			counter.Add(1)
			done <- struct{}{}
		})
		if err != nil {
			t.Fatalf("Submit: %v", err)
		}
	}

	for range n {
		<-done
	}
	if counter.Load() != n {
		t.Fatalf("counter = %d, want %d", counter.Load(), n)
	}
}

func TestPool_SubmitAfterStopReturnsError(t *testing.T) {
	p := New(2)
	p.Stop()

	err := p.Submit(context.Background(), func(context.Context) {})
	if !errors.Is(err, ErrPoolStopped) {
		t.Fatalf("err = %v, want ErrPoolStopped", err)
	}
}

func TestPool_SubmitRespectsContext(t *testing.T) {
	p := New(1)
	defer p.Stop()

	// Occupy the single worker with a blocking task.
	block := make(chan struct{})
	_ = p.Submit(context.Background(), func(context.Context) { <-block })

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	err := p.Submit(ctx, func(context.Context) {})
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("err = %v, want context.DeadlineExceeded", err)
	}
	close(block)
}

func TestPool_StopIsIdempotent(t *testing.T) {
	p := New(2)
	p.Stop()
	p.Stop() // must not panic
}
