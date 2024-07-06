package domain

import (
	"github.com/robfig/cron/v3"
	"time"
)

type CronJob struct {
	Id int64

	// 通用的任务的抽象
	// 具体任务设置具体的值
	Cfg string

	// 执行器名称
	Executor string

	// job 名称
	Name string

	// cron 表达式
	Expression string

	NextTime time.Time

	// 提供释放方法
	CancelFunc func()
}

var parser = cron.NewParser(cron.Second | cron.Minute | cron.Hour |
	cron.Dom | cron.Month | cron.Dow | cron.Descriptor)

// Next 下一次执行时间
func (j CronJob) Next(t time.Time) time.Time {
	schedule, _ := parser.Parse(j.Expression)
	return schedule.Next(t)
}
