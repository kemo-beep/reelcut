package worker

import (
	"context"
	"fmt"
	"time"

	"reelcut/internal/notifier"
	"reelcut/internal/queue"
	"reelcut/internal/repository"
	"reelcut/internal/service"

	"github.com/hibiken/asynq"
)

type RenderingWorker struct {
	renderingSvc *service.RenderingService
	clipRepo     repository.ClipRepository
	jobRepo      repository.ProcessingJobRepository
	notifier     notifier.JobNotifier
}

func NewRenderingWorker(renderingSvc *service.RenderingService, clipRepo repository.ClipRepository, jobRepo repository.ProcessingJobRepository, jobNotifier notifier.JobNotifier) *RenderingWorker {
	return &RenderingWorker{renderingSvc: renderingSvc, clipRepo: clipRepo, jobRepo: jobRepo, notifier: jobNotifier}
}

func (w *RenderingWorker) Register(mux *asynq.ServeMux) {
	mux.Handle(queue.TypeRender, asynq.HandlerFunc(w.Handle))
}

func (w *RenderingWorker) Handle(ctx context.Context, t *asynq.Task) error {
	payload, err := queue.ParseRenderPayload(t.Payload())
	if err != nil {
		return err
	}
	job, err := w.jobRepo.GetByID(ctx, payload.JobID)
	if err != nil || job == nil {
		return fmt.Errorf("job not found: %s", payload.JobID)
	}
	if job.Status == "cancelled" {
		return nil
	}
	job.Status = "processing"
	job.Progress = 10
	now := time.Now()
	job.StartedAt = &now
	_ = w.jobRepo.Update(ctx, job)
	if w.notifier != nil {
		w.notifier.NotifyJob(ctx, job)
	}

	if err := w.renderingSvc.Render(ctx, payload.ClipID, payload.Preset); err != nil {
		job.Status = "failed"
		if errMsg := err.Error(); errMsg != "" {
			job.ErrorMessage = &errMsg
		}
		_ = w.jobRepo.Update(ctx, job)
		if w.notifier != nil {
			w.notifier.NotifyJob(ctx, job)
		}
		c, _ := w.clipRepo.GetByID(ctx, payload.ClipID)
		if c != nil {
			c.Status = "failed"
			_ = w.clipRepo.Update(ctx, c)
		}
		return err
	}

	job.Progress = 100
	job.Status = "completed"
	completed := time.Now()
	job.CompletedAt = &completed
	if err := w.jobRepo.Update(ctx, job); err != nil {
		return err
	}
	if w.notifier != nil {
		w.notifier.NotifyJob(ctx, job)
	}
	return nil
}
