package client

import (
	"bytes"
	"context"
	"errors"
	"image/jpeg"
	"io"
	"log"
	"mini_tiktok/internal/dao/models"
	"mini_tiktok/internal/dao/mysql"
	publishService "mini_tiktok/internal/rpc/rpcGen/publish"
	service "mini_tiktok/internal/services"
	"mini_tiktok/utils"
	"mini_tiktok/utils/qiniu"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/bakape/thumbnailer/v2"
)

var PublishClient publishService.PublishServiceClient

const auth_address = "127.0.0.1:8891"

func init() {
	// 连接服务器
	conn, err := grpc.Dial(auth_address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("net.Connect err: %v", err)
	}
	defer conn.Close()

	// 建立gRPC连接
	PublishClient = publishService.NewPublishServiceClient(conn)
}

type PublishService struct {
	p *publishService.UnimplementedPublishServiceServer
}

func (p *PublishService) PublishAction(ctx context.Context, req *publishService.PublishActionRequest) (*publishService.PublishActionResponse, error) {
	filename := utils.GetFileName(req.UserId)

	// 视频上传七牛云对象存储并返回播放地址
	videoUrl, err := qiniu.UploadimageToQiNiu(req.Data, filename)
	if err != nil {
		return publishVideosError(err.Error()), err
	}

	// 生成缩略图并上传图片
	reader := bytes.NewReader(req.Data)
	thumbData, err := getThumbnail(reader)
	if err != nil {
		return publishVideosError(err.Error()), err
	}
	coverimageUrl, err := qiniu.UploadimageToQiNiu(thumbData, filename)
	if err != nil {
		return publishVideosError(err.Error()), err
	}
	if err := service.PostVideo(req.UserId, videoUrl, coverimageUrl, req.Title); err != nil {
		return publishVideosError(err.Error()), err
	}
	return publishVideosSuccess(), nil

}

func (p *PublishService) QueryPublishList(ctx context.Context, req *publishService.QueryPublishListRequest) (*publishService.QueryPublishListResponse, error) {
	videoList, err := service.PublishListService(req.UserId)
	if err != nil {
		return &publishService.QueryPublishListResponse{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		}, err
	}

	videos := []*publishService.Video{}
	for _, v := range *videoList {
		video, err := mongoVdoToPublishVdo(v, req.UserId)
		if err != nil {
			return &publishService.QueryPublishListResponse{
				StatusCode: -1,
				StatusMsg:  err.Error(),
			}, err
		}
		videos = append(videos, video)
	}
	return &publishService.QueryPublishListResponse{
		StatusCode: 0,
		StatusMsg:  "success",
		VideoList:  videos,
	}, nil
}

// getThumbnail Generate JPEG thumbnail from video
func getThumbnail(input io.ReadSeeker) ([]byte, error) {
	_, thumb, err := thumbnailer.Process(input, thumbnailer.Options{})
	if err != nil {
		return nil, errors.New("failed to create thumbnail")
	}
	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, thumb, nil)
	if err != nil {
		return nil, errors.New("failed to create buffer")
	}
	return buf.Bytes(), nil
}

func publishVideosError(msg string) *publishService.PublishActionResponse {
	return &publishService.PublishActionResponse{
		StatusCode: -1,
		StatusMsg:  msg,
	}
}

func publishVideosSuccess() *publishService.PublishActionResponse {
	return &publishService.PublishActionResponse{
		StatusCode: 0,
		StatusMsg:  "success",
	}
}

//将video.go中的Video转化为publish.pb.go中的video类型
func mongoVdoToPublishVdo(vdo *models.Video, userId int64) (*publishService.Video, error) {
	res := &publishService.Video{}

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
