package ratelimiter

import (
	"context"
	"testing"
	"time"
)

func TestWaitReturnsImmediatelyWithinBurst(t *testing.T) {
	l := New(1000) // large burst, generous rate: shouldn't block in practice
	defer l.Stop()

	ctx := context.Background()
	start := time.Now()
	for range 5 {
		if err := l.Wait(ctx); err != nil {
			t.Fatalf("Wait() unexpected error: %v", err)
		}
	}
	if elapsed := time.Since(start); elapsed > time.Second {
		t.Errorf("Wait() took %v for 5 tokens at rate=1000/s, expected well under 1s", elapsed)
	}
}

func TestWaitRespectsContextCancellation(t *testing.T) {
	// Rate of 1/s with an empty bucket means the second Wait has to queue;
	// cancel the context immediately and confirm Wait returns promptly
	// instead of blocking for the full interval.
	l := New(1)
	defer l.Stop()

	ctx := context.Background()
	if err := l.Wait(ctx); err != nil {
		t.Fatalf("first Wait() unexpected error: %v", err)
	}

	cancelCtx, cancel := context.WithCancel(context.Background())
	cancel()

	start := time.Now()
	err := l.Wait(cancelCtx)
	elapsed := time.Since(start)

	if err == nil {
		t.Fatal("Wait() expected error from cancelled context, got nil")
	}
	if elapsed > 100*time.Millisecond {
		t.Errorf("Wait() took %v to observe cancellation, expected near-instant", elapsed)
	}
}

func TestNewClampsNonPositiveRate(t *testing.T) {
	l := New(0)
	defer l.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := l.Wait(ctx); err != nil {
		t.Fatalf("Wait() with clamped rate should succeed within timeout, got: %v", err)
	}
}
