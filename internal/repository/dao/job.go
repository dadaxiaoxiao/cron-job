package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type GORMJobDAO struct {
	db *gorm.DB
}

func (g *GORMJobDAO) Stop(ctx context.Context, id int64) error {
	return g.db.WithContext(ctx).
		Model(&Job{}).
		Where("id= ?", id).
		Updates(map[string]any{
			"status": jobStatusPaused,
			"utime":  time.Now().UnixMilli(),
		}).Error
}

func (g *GORMJobDAO) UpdateNextTime(ctx context.Context, id int64, next time.Time) error {
	return g.db.WithContext(ctx).
		Model(&Job{}).
		Where("id= ?", id).
		Updates(map[string]any{
			"next_time": next.UnixMilli(),
			"utime":     time.Now().UnixMilli(),
		}).Error
}

// UpdateUTime 更新 UTime
func (g *GORMJobDAO) UpdateUTime(ctx context.Context, id int64) error {
	return g.db.WithContext(ctx).
		Model(&Job{}).
		Where("id= ?", id).
		Updates(map[string]any{
			"utime": time.Now().UnixMilli(),
		}).Error
}

// Release 释放为等待
func (g *GORMJobDAO) Release(ctx context.Context, id int64) error {
	return g.db.WithContext(ctx).
		Model(&Job{}).Where("id= ?", id).
		Updates(map[string]any{
			"status": jobStatusWaiting,
			"utime":  time.Now().UnixMilli(),
		}).Error
}

func (g *GORMJobDAO) Preempt(ctx context.Context) (Job, error) {
	// 高并发情况下
	// 假如有 100 个goroutine
	// 所有 goroutine 执行的循环次数加在一起是 100! = 1+2+3+....+100
	//  意味着特定一个 goroutine，最差情况下，要循环一百次
	db := g.db.WithContext(ctx)
	for {
		now := time.Now()
		var j Job
		// 分布式任务调度系统 高并发情况下
		// 1.一次拉一批，比如一次性取出 100 条来，然后，随机从某一条开始，向后开始抢占
		// 2. 搞个随机偏移量，0-100 生成一个随机偏移量。兜底：第一轮没查到，偏移量回归到 0
		// 3. 搞一个 id 取余分配，status = ? AND next_time <=? AND id%10 = ? 兜底：不加余数条件，取next_time 最早的
		err := db.Where("(status= ? AND next_time <= ?) OR (status= ? AND utime <= ?)", jobStatusWaiting, now, jobStatusRunning, now.Add(time.Minute*-3)).
			First(&j).Error

		if err != nil {
			// 没有任务
			return Job{}, err
		}
		// 高并发抢占任务
		// 使用乐观锁更新状态
		// utime CAS 操作，compare AND Swap
		// 就是用乐观锁取代 FOR UPDATE
		// 为什么不单纯使用 id ? 只使用id 每个 goroutine 都能成功更新，使用 utime CAS 操作，只有一个goroutine 操作成功
		res := db.Model(&Job{}).
			Where("id= ? AND version= ?", j.Id, j.Version).
			Updates(map[string]any{
				"status":  jobStatusRunning, // 标志抢占中
				"utime":   now,
				"version": j.Version + 1,
			})

		// 数据库相关错误
		if res.Error != nil {
			return Job{}, err
		}

		// 怎么判断抢占成功 ？
		if res.RowsAffected == 0 {
			// 抢占失败，继续下一轮
			continue
		}
		// 抢占成功
		return j, nil
	}
}

func (dao *GORMJobDAO) Insert(ctx context.Context, j Job) error {
	now := time.Now().UnixMilli()
	j.Ctime = now
	j.Utime = now
	return dao.db.WithContext(ctx).Create(&j).Error
}

func NewGORMJobDAO(db *gorm.DB) JobDAO {
	return &GORMJobDAO{
		db: db,
	}
}
