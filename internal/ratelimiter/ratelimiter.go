// Package ratelimiter provides a small token-bucket rate limiter used to
// bound how fast concurrent workers may proceed, independent of how many
// workers are running.
package ratelimiter

import (
	"context"
	"time"
)

// Limiter hands out tokens at a fixed rate, up to a burst of one second's
// worth of tokens. Multiple goroutines may call Wait concurrently.
type Limiter struct {
	tokens chan struct{}
	stop   chan struct{}
}

// New creates a Limiter that admits ratePerSecond tokens per second.
// Rates below 1 are clamped to 1.
func New(ratePerSecond int) *Limiter {
	if ratePerSecond < 1 {
		ratePerSecond = 1
	}

	l := &Limiter{
		tokens: make(chan struct{}, ratePerSecond), // burst capacity: 1s worth
		stop:   make(chan struct{}),
	}

	go l.refill(ratePerSecond)
	return l
}

func (l *Limiter) refill(ratePerSecond int) {
	interval := time.Second / time.Duration(ratePerSecond)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			select {
			case l.tokens <- struct{}{}:
			default:
				// Bucket is full; drop this tick rather than block the refiller.
			}
		case <-l.stop:
			return
		}
	}
}

// Wait blocks until a token is available or ctx is cancelled, whichever
// comes first.
func (l *Limiter) Wait(ctx context.Context) error {
	select {
	case <-l.tokens:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Stop releases the background refill goroutine. Safe to call once.
func (l *Limiter) Stop() {
	close(l.stop)
}
