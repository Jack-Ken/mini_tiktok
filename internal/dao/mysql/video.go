package mysql

import (
	"errors"
	"mini_tiktok/internal/dao/models"
	db "mini_tiktok/internal/initialize"
	"time"
)

// 获取视频
func QueryVideoListByLimtAndTime(limit int, latestTime time.Time, videoList *[]*models.Video) error {
	if videoList == nil {
		return errors.New("空指针错误")
	}
	return db.SqlSession.Model(&models.Video{}).Where("created_at<?", latestTime).Order("created_at ASC").
		Limit(limit).Find(videoList).Error
}
