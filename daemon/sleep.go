package daemon

import (
	"context"
	"time"
)

// Sleep with context, returns false if context was canceled
func Sleep(ctx context.Context, d time.Duration) bool {
	t := time.NewTimer(d)
	select {
	case <-ctx.Done():
		if !t.Stop() {
			<-t.C
		}
		return false
	case <-t.C:
		return true
	}
}
