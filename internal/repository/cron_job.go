package repository

import (
	"context"
	"github.com/dadaxiaoxiao/cron-job/internal/domain"
	"github.com/dadaxiaoxiao/cron-job/internal/repository/dao"
	"time"
)

type PreemptCronJobRepository struct {
	dao dao.JobDAO
}

func NewPreemptCronJobRepository(dao dao.JobDAO) CronJobRepository {
	return &PreemptCronJobRepository{
		dao: dao,
	}
}

func (p *PreemptCronJobRepository) Release(ctx context.Context, id int64) error {
	return p.dao.Release(ctx, id)
}

func (p *PreemptCronJobRepository) Preempt(ctx context.Context) (domain.CronJob, error) {
	j, err := p.dao.Preempt(ctx)
	if err != nil {
		return domain.CronJob{}, err
	}
	return p.entityToDomain(j), nil
}

func (p *PreemptCronJobRepository) Stop(ctx context.Context, id int64) error {
	return p.dao.Stop(ctx, id)
}

func (p *PreemptCronJobRepository) UpdateUTime(ctx context.Context, id int64) error {
	return p.dao.UpdateUTime(ctx, id)
}

func (p *PreemptCronJobRepository) UpdateNextTime(ctx context.Context, id int64, next time.Time) error {
	return p.dao.UpdateNextTime(ctx, id, next)
}

func (p *PreemptCronJobRepository) AddJob(ctx context.Context, j domain.CronJob) error {
	return p.dao.Insert(ctx, p.toEntity(j))
}

func (p *PreemptCronJobRepository) entityToDomain(j dao.Job) domain.CronJob {
	return domain.CronJob{
		Id:         j.Id,
		Cfg:        j.Cfg,
		Executor:   j.Executor,
		Name:       j.Name,
		Expression: j.Expression,
		NextTime:   time.UnixMilli(j.NextTime),
	}
}

func (p *PreemptCronJobRepository) toEntity(j domain.CronJob) dao.Job {
	return dao.Job{
		Id:         j.Id,
		Name:       j.Name,
		Expression: j.Expression,
		Cfg:        j.Cfg,
		Executor:   j.Executor,
		NextTime:   j.NextTime.UnixMilli(),
	}
}
