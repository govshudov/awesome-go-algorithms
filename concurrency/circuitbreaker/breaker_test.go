package circuitbreaker

import (
	"errors"
	"testing"
	"time"
)

var errBoom = errors.New("boom")

func TestBreaker_OpensAfterThreshold(t *testing.T) {
	cb := New(Config{FailureThreshold: 3, ResetTimeout: time.Hour})
	for range 3 {
		if err := cb.Do(func() error { return errBoom }); err != errBoom {
			t.Fatalf("err = %v, want errBoom", err)
		}
	}
	if cb.State() != StateOpen {
		t.Fatalf("state = %s, want open", cb.State())
	}
	// Now the breaker short-circuits.
	if err := cb.Do(func() error { return nil }); err != ErrOpen {
		t.Fatalf("err = %v, want ErrOpen", err)
	}
}

func TestBreaker_SuccessResetsFailureCount(t *testing.T) {
	cb := New(Config{FailureThreshold: 3, ResetTimeout: time.Hour})
	_ = cb.Do(func() error { return errBoom })
	_ = cb.Do(func() error { return errBoom })
	// One success resets the counter.
	_ = cb.Do(func() error { return nil })
	// Two more failures should not trip (threshold is 3 consecutive).
	_ = cb.Do(func() error { return errBoom })
	_ = cb.Do(func() error { return errBoom })
	if cb.State() != StateClosed {
		t.Fatalf("state = %s, want closed", cb.State())
	}
}

func TestBreaker_HalfOpenRecovery(t *testing.T) {
	cb := New(Config{FailureThreshold: 1, ResetTimeout: 10 * time.Millisecond})
	virtual := time.Now()
	cb.now = func() time.Time { return virtual }

	_ = cb.Do(func() error { return errBoom })
	if cb.state != StateOpen {
		t.Fatalf("expected open after first failure, got %s", cb.state)
	}

	// Advance past reset timeout; the next call should probe in half-open.
	virtual = virtual.Add(20 * time.Millisecond)

	if err := cb.Do(func() error { return nil }); err != nil {
		t.Fatalf("probe err = %v", err)
	}
	if cb.State() != StateClosed {
		t.Fatalf("state = %s, want closed after successful probe", cb.State())
	}
}

func TestBreaker_HalfOpenFailureReopens(t *testing.T) {
	cb := New(Config{FailureThreshold: 1, ResetTimeout: 10 * time.Millisecond})
	virtual := time.Now()
	cb.now = func() time.Time { return virtual }

	_ = cb.Do(func() error { return errBoom })
	virtual = virtual.Add(20 * time.Millisecond)

	if err := cb.Do(func() error { return errBoom }); err != errBoom {
		t.Fatalf("probe err = %v, want errBoom", err)
	}
	if cb.State() != StateOpen {
		t.Fatalf("state = %s, want open after failed probe", cb.State())
	}
}
