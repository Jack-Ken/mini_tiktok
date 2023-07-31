package client

import (
	"context"
	"errors"
	"mini_tiktok/internal/dao/models"
	"mini_tiktok/internal/dao/mysql"
	feedService "mini_tiktok/internal/rpc/rpcGen/feed"
	service "mini_tiktok/internal/services"
	"mini_tiktok/utils"
	"strconv"
	"time"
)

type FeedService struct {
	*feedService.UnimplementedFeedServiceServer
}

func (feed *FeedService) Feed(ctx context.Context, req *feedService.FeedRequest) (*feedService.FeedResponse, error) {
	p := NewProxyFeedVideoList(req)
	var err error
	// Token为空

	if p.Req.Token == "" {
		err = p.DoNoToken()
	} else {
		err = p.DoWithToken()
	}

	if err != nil {
		return p.Resp, err
	}
	return p.Resp, nil
}

type ProxyFeedVideoList struct {
	Req  *feedService.FeedRequest
	Resp *feedService.FeedResponse
	uid  int64
}

func NewProxyFeedVideoList(req *feedService.FeedRequest) *ProxyFeedVideoList {
	return &ProxyFeedVideoList{
		Req:  req,
		Resp: &feedService.FeedResponse{},
	}
}

func (p *ProxyFeedVideoList) DoWithToken() error {
	mc, err := utils.ParseToken(p.Req.Token)
	if err != nil {
		return errors.New("token不正确")
	}
	if time.Now().Unix() > mc.ExpiresAt.Unix() {
		return errors.New("token超时")
	}
	p.uid = mc.UserID

	rawTimeStamp := p.Req.LatestTime
	var latestTime time.Time
	if rawTimeStamp == "" {
		latestTime = time.Unix(time.Now().Unix(), 0)
	} else {
		intTime, err := strconv.ParseInt(rawTimeStamp, 10, 64)
		if err != nil {
			latestTime = time.Unix(0, intTime*1e6) //注意：前端传来的时间戳是以ms为单位的
		}
	}

	// 调用service层接口，获取videoList
	videoList, err := service.QueryFeedVideoList(p.uid, latestTime)

	if err != nil {
		p.FeedVideoListError(err.Error())
		return err
	}

	videos := []*feedService.Video{}
	for i := 0; i < len(videoList.VideoList); i++ {
		video, err := mongoVdoToBizVdo(videoList.VideoList[i], p.uid)
		if err != nil {
			p.FeedVideoListError(err.Error())
			return err
		}
		videos = append(videos, video)
	}

	p.Resp.VideoList = videos
	p.Resp.NextTime = videoList.NextTime

	p.FeedVideoListSuccess()
	// 调用service层接口，获取videoList
	return nil
}

func (p *ProxyFeedVideoList) DoNoToken() error {
	rawTimeStamp := p.Req.LatestTime
	var latestTime time.Time
	if rawTimeStamp == "" {
		latestTime = time.Unix(time.Now().Unix(), 0)
	} else {
		intTime, err := strconv.ParseInt(rawTimeStamp, 10, 64)
		if err != nil {
			latestTime = time.Unix(0, intTime*1e6) //注意：前端传来的时间戳是以ms为单位的
		}
	}

	// 调用service层接口，获取videoList
	videoList, err := service.QueryFeedVideoList(0, latestTime)
	if err != nil {
		p.FeedVideoListError(err.Error())
		return err
	}

	videos := []*feedService.Video{}
	for i := 0; i < len(videoList.VideoList); i++ {
		video, err := mongoVdoToBizVdo(videoList.VideoList[i], p.uid)
		if err != nil {
			p.FeedVideoListError(err.Error())
			return err
		}
		videos = append(videos, video)
	}

	p.Resp.VideoList = videos
	p.Resp.NextTime = videoList.NextTime

	p.FeedVideoListSuccess()
	return nil
}

func (p *ProxyFeedVideoList) FeedVideoListError(msg string) {
	p.Resp = &feedService.FeedResponse{
		StatusCode: -1,
		StatusMsg:  msg,
	}
}

func (p *ProxyFeedVideoList) FeedVideoListSuccess() {
	p.Resp.StatusMsg = "success"
	p.Resp.StatusCode = 0

}

//将video.go中的Video转化为feed.pb.go中的video类型
func mongoVdoToBizVdo(vdo *models.Video, userId int64) (*feedService.Video, error) {
	res := &feedService.Video{}

	res.Id = vdo.Id

	// 获取视频作者信息
	videoUser, err := service.Info_Service(vdo.UserId)
	if err != nil {
		return res, err
	}

	res.Author.FollowerCount = videoUser.FollowerCount
	res.Author.Id = videoUser.Id
	res.Author.FollowCount = videoUser.FollowCount
	res.Author.IsFollow = videoUser.IsFollow
	res.Author.Name = videoUser.Name

	res.PlayUrl = vdo.PlayUrl
	res.CoverUrl = vdo.CoverUrl
	res.FavoriteCount = vdo.FavoriteCount
	res.CommentCount = int64(len(vdo.Comments))
	res.Title = vdo.Title

	//判断当前用户是否点赞
	f1, _ := mysql.GetVideoFavorState(userId, vdo.Id)

	res.IsFavorite = f1

	return res, nil
}
