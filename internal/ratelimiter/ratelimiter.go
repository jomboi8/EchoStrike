package ratelimiter

import (
	"context"
	"time"
)

type Limiter struct {
	tokens chan struct{}
	stop   chan struct{}
}

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
