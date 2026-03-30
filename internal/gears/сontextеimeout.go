package gears

import (
	"context"
	"time"
)

func ContextTimeout(ctx context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	if timeout > 0 {
		return context.WithTimeout(ctx, timeout)
	}

	return ctx, func() {}
}
