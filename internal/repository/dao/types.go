package dao

import (
	"context"
	"time"
)

type JobDAO interface {
	// Preempt 抢占
	Preempt(ctx context.Context) (Job, error)
	// Release 释放
	Release(ctx context.Context, id int64) error
	// Stop 停止
	Stop(ctx context.Context, id int64) error
	UpdateUTime(ctx context.Context, id int64) error
	UpdateNextTime(ctx context.Context, id int64, next time.Time) error
	Insert(ctx context.Context, j Job) error
}

type Job struct {
	Id int64 `gorm:"primaryKey,autoIncrement"`

	// 通用的任务的抽象
	Cfg string `gorm:"column:cfg;type:text;comment:通用的任务的抽象"`

	// 执行器名称
	Executor string `gorm:"column:executor;type:varchar(255);comment:执行器名称"`

	// 任务名称
	Name string `gorm:"unique;column:name;type:varchar(128);comment:任务名称"`

	// 状态 jobStatusWaiting ，jobStatusRunning，jobStatusPaused
	Status int `gorm:"column:status;comment:状态"`

	// 下一次被调度的时间
	NextTime int64 `gorm:"index;column:next_time;comment:下一次被调度的时间"`

	// cron 表达式
	Expression string `gorm:"column:expression;type:varchar(128);comment:cron表达式"`

	Version int `gorm:"column:version;comment:版本号"`

	// 创建时间 毫秒数
	Ctime int64

	// 更新时间 毫秒数
	Utime int64
}

const (
	// 等待
	jobStatusWaiting = iota
	// 已经被抢占进行
	jobStatusRunning
	// 暂停
	jobStatusPaused
)
