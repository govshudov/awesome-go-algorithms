package ratelimiter

import (
	"context"
	"testing"
	"time"
)

func TestTokenBucket_InitialCapacityConsumed(t *testing.T) {
	b := New(3, time.Second)
	// Pin virtual clock so refill does not add tokens mid-test.
	start := time.Now()
	b.now = func() time.Time { return start }

	for i := range 3 {
		if !b.Allow() {
			t.Fatalf("Allow #%d = false, want true", i)
		}
	}
	if b.Allow() {
		t.Fatal("Allow after exhausting capacity = true, want false")
	}
}

func TestTokenBucket_Refills(t *testing.T) {
	b := New(2, 10*time.Millisecond)
	virtual := time.Now()
	b.now = func() time.Time { return virtual }

	b.Allow()
	b.Allow()
	if b.Allow() {
		t.Fatal("bucket should be empty")
	}

	// Advance virtual clock 25ms -> 2 refills queued, capped at capacity=2.
	virtual = virtual.Add(25 * time.Millisecond)
	if !b.Allow() {
		t.Fatal("expected refilled token")
	}
	if !b.Allow() {
		t.Fatal("expected second refilled token")
	}
	if b.Allow() {
		t.Fatal("refill exceeded capacity")
	}
}

func TestTokenBucket_WaitRespectsContext(t *testing.T) {
	b := New(1, 10*time.Second)
	b.Allow() // drain

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	err := b.Wait(ctx)
	if err != context.DeadlineExceeded {
		t.Fatalf("err = %v, want DeadlineExceeded", err)
	}
}

func TestTokenBucket_WaitSucceeds(t *testing.T) {
	b := New(1, 5*time.Millisecond)
	b.Allow() // drain
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	if err := b.Wait(ctx); err != nil {
		t.Fatalf("Wait err = %v", err)
	}
}
