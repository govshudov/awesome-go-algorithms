// Package workerpool provides a fixed-size worker pool with graceful
// shutdown and per-task context cancellation.
package workerpool

import (
	"context"
	"errors"
	"sync"
)

// ErrPoolStopped is returned by Submit after Stop has been called.
var ErrPoolStopped = errors.New("workerpool: stopped")

type Task func(ctx context.Context)

type Pool struct {
	tasks  chan Task
	wg     sync.WaitGroup
	stop   chan struct{}
	once   sync.Once
	stopMu sync.RWMutex
	closed bool
}

func New(workers int) *Pool {
	if workers < 1 {
		workers = 1
	}
	p := &Pool{
		tasks: make(chan Task),
		stop:  make(chan struct{}),
	}
	p.wg.Add(workers)
	for range workers {
		go p.worker()
	}
	return p
}

func (p *Pool) worker() {
	defer p.wg.Done()
	for task := range p.tasks {
		task(p.ctx())
	}
}

// Submit enqueues a task. It returns ErrPoolStopped if the pool has been
// stopped, or ctx.Err() if the caller's context is cancelled while waiting.
func (p *Pool) Submit(ctx context.Context, task Task) error {
	p.stopMu.RLock()
	if p.closed {
		p.stopMu.RUnlock()
		return ErrPoolStopped
	}
	p.stopMu.RUnlock()

	select {
	case p.tasks <- task:
		return nil
	case <-p.stop:
		return ErrPoolStopped
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Stop closes the task channel and waits for in-flight tasks to finish.
// Safe to call multiple times.
func (p *Pool) Stop() {
	p.once.Do(func() {
		p.stopMu.Lock()
		p.closed = true
		p.stopMu.Unlock()
		close(p.stop)
		close(p.tasks)
	})
	p.wg.Wait()
}

// ctx returns a context that is cancelled when the pool is stopped.
func (p *Pool) ctx() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-p.stop
		cancel()
	}()
	return ctx
}
