// Package timer - Simple safety time sleep.
package timer

import (
	"context"
	"time"
)

// Run - Simple safety time sleep.
func Run(d time.Duration) error {
	return RunWithContext(context.Background(), d)
}

// RunWithContext - Simple safety time sleep.
// Priority is given to canceling the context.
func RunWithContext(ctx context.Context, d time.Duration) error {
	t := time.NewTimer(d)
	defer t.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.C:
		break
	}

	return nil
}
