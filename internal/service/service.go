package service

import (
	"context"
	"github.com/dadaxiaoxiao/cron-job/internal/domain"
	"github.com/dadaxiaoxiao/cron-job/internal/repository"
	"github.com/dadaxiaoxiao/go-pkg/accesslog"
	"time"
)

type CronJobService interface {
	Preempt(ctx context.Context) (domain.CronJob, error)
	ResetNextTime(ctx context.Context, job domain.CronJob) error
	AddJob(ctx context.Context, j domain.CronJob) error
	Stop(ctx context.Context, job domain.CronJob) error
}

type cronJobService struct {
	// 续约机制
	refreshInterval time.Duration
	repo            repository.CronJobRepository
	log             accesslog.Logger
}

func NewCronJobService(repo repository.CronJobRepository, log accesslog.Logger) CronJobService {
	return &cronJobService{
		repo:            repo,
		refreshInterval: time.Second * 10,
		log:             log,
	}
}

// Preempt 抢占
func (c *cronJobService) Preempt(ctx context.Context) (domain.CronJob, error) {
	job, err := c.repo.Preempt(ctx)
	if err != nil {
		return domain.CronJob{}, err
	}

	ch := make(chan struct{})
	go func() {
		// 续约
		// 启动goroutine 开始续约，也就是持续抢占，
		// 这里是每隔着 refreshInterval 进行一次续约
		ticker := time.NewTicker(c.refreshInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ch:
				// 退出续约循环
				return
			case <-ticker.C:
				c.refresh(job.Id)
			}
		}
	}()

	// 抢占后，是否要一直抢占？
	// 提供释放
	job.CancelFunc = func() {
		close(ch)
		// 这里新建一个ctx 来控制超时
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		err := c.repo.Release(ctx, job.Id)
		if err != nil {
			c.log.Error("释放任务失败",
				accesslog.Error(err),
				accesslog.Int64("id", job.Id))
		}
	}
	return job, nil
}

func (c *cronJobService) ResetNextTime(ctx context.Context, job domain.CronJob) error {
	// 计算下一次的时间
	t := job.Next(time.Now())
	if !t.IsZero() {
		return c.repo.UpdateNextTime(ctx, job.Id, t)
	}
	return nil
}

func (c *cronJobService) AddJob(ctx context.Context, j domain.CronJob) error {
	j.NextTime = j.Next(time.Now())
	return c.repo.AddJob(ctx, j)
}

func (c *cronJobService) Stop(ctx context.Context, j domain.CronJob) error {
	return c.repo.Stop(ctx, j.Id)
}

// refresh 刷新续约
func (c *cronJobService) refresh(id int64) {
	// 这里新建一个ctx 来控制超时
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	// 续约逻辑
	// 更新job 的更新时间即好
	// 续约成功 : 处于 running 状态，但是更新时间在三分钟以内
	// 续约失败 : 处于 running 状态，但更新时间在三分钟之前 （可以其他实例抢占）
	err := c.repo.UpdateUTime(ctx, id)
	if err != nil {
		c.log.Error("job续约失败",
			accesslog.Error(err),
			accesslog.Int64("jobId", id))
	}
}
