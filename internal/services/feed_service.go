package services

import (
	"mini_tiktok/internal/dao/models"
	"mini_tiktok/internal/dao/mysql"
	"time"
)

const MaxVideoNum = 30

type FeedVideoList struct {
	VideoList []*models.Video `json:"video_list,omitempty"`
	NextTime  int64           `json:"next_time,omitempty"`
}

type QueryFeed struct {
	userid     int64
	latestTime time.Time
	videos     []*models.Video
	nextTime   int64
	feedVideo  *FeedVideoList
}

func QueryFeedVideoList(uid int64, latestTime time.Time) (*FeedVideoList, error) {
	return NewQueryFeedVideoList(uid, latestTime).Do()
}

func NewQueryFeedVideoList(uid int64, latestTime time.Time) *QueryFeed {
	return &QueryFeed{userid: uid, latestTime: latestTime}
}

func (q *QueryFeed) Do() (*FeedVideoList, error) {
	q.checkParams()
	if err := q.prepareData(); err != nil {
		return nil, err
	}
	q.feedVideo = &FeedVideoList{
		VideoList: q.videos,
		NextTime:  q.nextTime,
	}
	//
	return q.feedVideo, nil

}

func (q *QueryFeed) checkParams() {
	if q.userid > 0 {
		return
		//uid有效，可以做一次定制的推荐视频流
	}
	if q.latestTime.IsZero() {
		q.latestTime = time.Now()
	}
}
func (q *QueryFeed) prepareData() error {
	err := mysql.QueryVideoListByLimtAndTime(MaxVideoNum, q.latestTime, &q.videos)
	if err != nil {
		return err
	}
	// 准备下一次调用的时间戳
	q.nextTime = time.Now().Unix() / 1e6
	return nil

}
