package controller

import (
	"context"
	"log"
	feedService "mini_tiktok/internal/rpc/rpcGen/feed"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var FeedClient feedService.FeedServiceClient

const feed_address = "127.0.0.1:8890"

func init() {
	// 连接服务器
	conn, err := grpc.Dial(feed_address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("net.Connect err: %v", err)
	}
	defer conn.Close()

	// 建立gRPC连接
	FeedClient = feedService.NewFeedServiceClient(conn)
}

func Feed_Hanlder(c *gin.Context) {
	token := c.Query("token")
	rawTimeStamp := c.Query("latest_time")

	resp, err := FeedClient.Feed(context.Background(), &feedService.FeedRequest{
		Token:      token,
		LatestTime: rawTimeStamp,
	})
	if err != nil {
		zap.L().Error("get video feed list failed", zap.Error(err))
	}
	c.JSON(http.StatusOK, resp)
}

//type ProxyFeedVideoList struct {
//	*gin.Context
//}
//
//func NewProxyFeedVideoList(c *gin.Context) *ProxyFeedVideoList {
//	return &ProxyFeedVideoList{c}
//}
//
//// 未登录的视频流处理
//func (p *ProxyFeedVideoList) DoNoToken() error {
//	rawTimeStamp := p.Query("latest_time")
//	var latestTime time.Time
//	if rawTimeStamp == "" {
//		latestTime = time.Unix(time.Now().Unix(), 0)
//	} else {
//		intTime, err := strconv.ParseInt(rawTimeStamp, 10, 64)
//		if err != nil {
//			latestTime = time.Unix(0, intTime*1e6) //注意：前端传来的时间戳是以ms为单位的
//		}
//	}
//	//intTime, err := strconv.ParseInt(rawTimeStamp, 10, 64)
//	//if err != nil {
//	//	latestTime = time.Unix(0, intTime*1e6) //注意：前端传来的时间戳是以ms为单位的
//	//}
//
//	// 调用service层接口，获取videoList
//	videoList, err := service.QueryFeedVideoList(0, latestTime)
//	if err != nil {
//		return err
//	}
//	p.FeedVideoListSuccess(videoList)
//	return nil
//}
//
//// 登录的视频流处理
//func (p *ProxyFeedVideoList) DoWithToken(token string) error {
//	mc, err := utils.ParseToken(token)
//	if err != nil {
//		return errors.New("token不正确")
//	}
//	if time.Now().Unix() > mc.ExpiresAt.Unix() {
//		return errors.New("token超时")
//	}
//	rawTimeStamp := p.Query("latest_time")
//	var latestTime time.Time
//	if rawTimeStamp == "" {
//		latestTime = time.Unix(time.Now().Unix(), 0)
//	} else {
//		intTime, err := strconv.ParseInt(rawTimeStamp, 10, 64)
//		if err != nil {
//			latestTime = time.Unix(0, intTime*1e6) //注意：前端传来的时间戳是以ms为单位的
//		}
//	}
//	//fmt.Println("latestTime is:", latestTime)
//	// 调用service层接口，获取videoList
//	videoList, err := service.QueryFeedVideoList(mc.UserID, latestTime)
//	if err != nil {
//		return err
//	}
//	p.FeedVideoListSuccess(videoList)
//	// 调用service层接口，获取videoList
//	return nil
//
//}
//
//func (p *ProxyFeedVideoList) FeedVideoListError(msg string) {
//	p.JSON(http.StatusOK, models.FeedResponse{
//		Response: models.Response{
//			StatusCode: 1,
//			StatusMsg:  msg,
//		},
//	})
//}
//
//func (p *ProxyFeedVideoList) FeedVideoListSuccess(videoList *service.FeedVideoList) {
//	p.JSON(http.StatusOK, models.FeedResponse{
//		Response: models.Response{
//			StatusCode: 0,
//			StatusMsg:  "success",
//		},
//		VideoList: videoList.VideoList,
//		NextTime:  videoList.NextTime,
//	})
//}
