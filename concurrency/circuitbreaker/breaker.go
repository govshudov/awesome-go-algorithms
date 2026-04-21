// Package circuitbreaker implements a classic three-state circuit breaker
// (Closed / Open / HalfOpen) that short-circuits calls to a failing
// dependency and probes for recovery after a reset timeout.
package circuitbreaker

import (
	"errors"
	"sync"
	"time"
)

// State of the breaker.
type State int

const (
	StateClosed State = iota
	StateOpen
	StateHalfOpen
)

func (s State) String() string {
	switch s {
	case StateClosed:
		return "closed"
	case StateOpen:
		return "open"
	case StateHalfOpen:
		return "half-open"
	default:
		return "unknown"
	}
}

// ErrOpen is returned by Do when the breaker is open.
var ErrOpen = errors.New("circuitbreaker: open")

type Config struct {
	// FailureThreshold is the number of consecutive failures required to
	// trip the breaker open. Must be >= 1.
	FailureThreshold int
	// ResetTimeout is how long to stay open before allowing a probe.
	ResetTimeout time.Duration
}

type Breaker struct {
	cfg      Config
	mu       sync.Mutex
	state    State
	failures int
	openedAt time.Time
	now      func() time.Time
}

func New(cfg Config) *Breaker {
	if cfg.FailureThreshold < 1 {
		cfg.FailureThreshold = 1
	}
	if cfg.ResetTimeout <= 0 {
		cfg.ResetTimeout = time.Second
	}
	return &Breaker{cfg: cfg, state: StateClosed, now: time.Now}
}

// State returns the current state, applying the reset-timeout transition if
// applicable.
func (b *Breaker) State() State {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.maybeHalfOpen()
	return b.state
}

// Do executes fn, recording the result and updating the breaker state.
// When open, Do returns ErrOpen without calling fn.
func (b *Breaker) Do(fn func() error) error {
	b.mu.Lock()
	b.maybeHalfOpen()
	if b.state == StateOpen {
		b.mu.Unlock()
		return ErrOpen
	}
	b.mu.Unlock()

	err := fn()

	b.mu.Lock()
	defer b.mu.Unlock()
	if err != nil {
		b.recordFailure()
		return err
	}
	b.recordSuccess()
	return nil
}

func (b *Breaker) maybeHalfOpen() {
	if b.state == StateOpen && b.now().Sub(b.openedAt) >= b.cfg.ResetTimeout {
		b.state = StateHalfOpen
	}
}

func (b *Breaker) recordFailure() {
	b.failures++
	if b.state == StateHalfOpen || b.failures >= b.cfg.FailureThreshold {
		b.state = StateOpen
		b.openedAt = b.now()
	}
}

func (b *Breaker) recordSuccess() {
	b.failures = 0
	b.state = StateClosed
}
