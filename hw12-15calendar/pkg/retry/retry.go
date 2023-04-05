package retry

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
)

type Connector func(ctx context.Context) error

// Retry returns a function matching the retryConnector type that
// is trying to establish a connection with the database retries number
// every delay time.
func Retry(connector Connector, retries int) Connector {
	return func(ctx context.Context) error {
		for r := 0; ; r++ {
			err := connector(ctx)
			if err == nil || r >= retries {
				return err
			}

			// Exponential increase in latency.
			shouldRetryAt := time.Second * 2 << r
			zap.S().Warnf("Attempt %d failed; retrying in %v", r+1, shouldRetryAt)

			select {
			case <-time.After(shouldRetryAt):
			case <-ctx.Done():
				return fmt.Errorf("retry: %w", ctx.Err())
			}
		}
	}
}
