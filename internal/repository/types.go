package repository

import (
	"context"
	"github.com/dadaxiaoxiao/cron-job/internal/domain"
	"time"
)

type CronJobRepository interface {
	Preempt(ctx context.Context) (domain.CronJob, error)
	Release(ctx context.Context, id int64) error
	Stop(ctx context.Context, id int64) error
	UpdateUTime(ctx context.Context, id int64) error
	UpdateNextTime(ctx context.Context, id int64, t time.Time) error
	AddJob(ctx context.Context, j domain.CronJob) error
}
