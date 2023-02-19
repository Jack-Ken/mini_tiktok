package services

import (
	"mini_tiktok/internal/dao/models"
	"mini_tiktok/internal/dao/mysql"
)

type QueryPublishListFlow struct {
	uid    int64
	videos []*models.Video
}

func PublishListService(uid int64) (*[]*models.Video, error) {
	return NewQueryPublishListFlow(uid).Do()
}

func NewQueryPublishListFlow(uid int64) *QueryPublishListFlow {
	return &QueryPublishListFlow{uid: uid}
}

func (q *QueryPublishListFlow) Do() (*[]*models.Video, error) {
	if err := q.checkUser(); err != nil {
		return nil, err
	}
	if err := q.prepareData(); err != nil {
		return nil, err
	}
	return &q.videos, nil

}

func (q *QueryPublishListFlow) checkUser() error {
	// 检查用户是否已经存在
	return mysql.CheckUserExitById(q.uid)
}

func (q *QueryPublishListFlow) prepareData() error {
	// 在video中添加user信息
	//注意：Video由于在数据库中没有存储作者信息，所以需要手动填充
	//获取用户视频信息
	if err := mysql.QueryVideoListByUserId(q.uid, &q.videos); err != nil {
		return err
	}
	//获取用户信息做信息填充
	user, err := mysql.UserInfo(q.uid)
	if err != nil {
		return err
	}

	//填充信息(Author和IsFavorite字段
	for i := range q.videos {
		q.videos[i].Author = *user
		isFavorite, errs := mysql.GetVideoFavorState(q.uid, q.videos[i].Id)
		if errs != nil {
			return err
		}
		q.videos[i].IsFavorite = isFavorite
	}
	return nil
}
