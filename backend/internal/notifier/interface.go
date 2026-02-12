package notifier

import (
	"context"

	"reelcut/internal/domain"
)

// JobNotifier sends job updates to connected clients (e.g. over WebSocket).
// If nil, workers simply do not notify.
type JobNotifier interface {
	NotifyJob(ctx context.Context, job *domain.ProcessingJob)
}
