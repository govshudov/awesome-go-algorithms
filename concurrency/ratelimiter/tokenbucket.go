// Package ratelimiter provides a goroutine-safe token-bucket rate limiter.
//
// The bucket refills at a configured rate up to a maximum capacity. Allow
// returns true if a token is immediately available; Wait blocks until one
// becomes available or ctx is cancelled.
package ratelimiter

import (
	"context"
	"sync"
	"time"
)

type TokenBucket struct {
	capacity   int
	refillRate time.Duration
	mu         sync.Mutex
	tokens     int
	lastRefill time.Time
	now        func() time.Time
}

// New creates a bucket that holds up to capacity tokens and refills one
// token every refillInterval.
func New(capacity int, refillInterval time.Duration) *TokenBucket {
	if capacity < 1 {
		capacity = 1
	}
	if refillInterval <= 0 {
		refillInterval = time.Millisecond
	}
	return &TokenBucket{
		capacity:   capacity,
		refillRate: refillInterval,
		tokens:     capacity,
		lastRefill: time.Now(),
		now:        time.Now,
	}
}

// Allow returns true and consumes a token if one is available. Non-blocking.
func (b *TokenBucket) Allow() bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.refill()
	if b.tokens > 0 {
		b.tokens--
		return true
	}
	return false
}

// Wait blocks until a token is available or ctx is cancelled.
func (b *TokenBucket) Wait(ctx context.Context) error {
	for {
		if b.Allow() {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(b.refillRate):
		}
	}
}

func (b *TokenBucket) refill() {
	now := b.now()
	elapsed := now.Sub(b.lastRefill)
	if elapsed < b.refillRate {
		return
	}
	add := int(elapsed / b.refillRate)
	b.tokens += add
	if b.tokens > b.capacity {
		b.tokens = b.capacity
	}
	b.lastRefill = b.lastRefill.Add(time.Duration(add) * b.refillRate)
}
