package controller

import (
	"context"
	"fmt"
	"mime/multipart"
	"mini_tiktok/internal/dao/models"
	"mini_tiktok/internal/rpc/client"
	"mini_tiktok/internal/rpc/rpcGen/publish"
	"net/http"
	"path/filepath"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

var (
	videoIndexMap = map[string]struct{}{
		".mp4":  {},
		".avi":  {},
		".wmv":  {},
		".flv":  {},
		".mpeg": {},
		".mov":  {},
	}
	pictureIndexMap = map[string]struct{}{
		".jpg": {},
		".bmp": {},
		".png": {},
		".svg": {},
	}
)

func PublishVideosHandler(c *gin.Context) {
	// 参数准备
	uid := c.GetInt64("userid")

	title := c.PostForm("title")
	fmt.Println(title)
	// 上传多文件
	form, err := c.MultipartForm()
	if err != nil {
		publishVideosError(c, err.Error())
	}
	files := form.File["data"]

	// 遍历所有文件
	for _, file := range files {
		opened, _ := file.Open()

		defer func(opened multipart.File) {
			err := opened.Close()
			if err != nil {
				zap.L().Error("opened.Close() failed", zap.Error(err))
			}
		}(opened)

		// file转byte类型
		var data = make([]byte, file.Size)
		readSize, err := opened.Read(data)
		if err != nil {
			zap.L().Error("file read to buye failed", zap.Error(err))
			publishVideosError(c, err.Error())
			return
		}

		if readSize != int(file.Size) {
			zap.L().Error("file read to buye failed", zap.Error(err))
			publishVideosError(c, err.Error())
			return
		}

		// 检查文件属性是否是视频
		suffix := filepath.Ext(file.Filename) //得到后缀
		if _, ok := videoIndexMap[suffix]; !ok {
			publishVideosError(c, "不支持的视频类型")
		}

		resp, err := client.PublishClient.PublishAction(context.Background(), &publish.PublishActionRequest{
			UserId: uid,
			Data:   data,
			Title:  title,
		})

		c.JSON(http.StatusOK, resp)

		//// 获取文件名称并存储
		//name := utils.GetFileName(uid)
		//filename := name + suffix
		//savePath := filepath.Join("../static/videos", filename)
		//if err := c.SaveUploadedFile(file, savePath); err != nil {
		//	publishVideosError(c, fmt.Sprintf("upload err %s", err.Error()))
		//	return
		//}
		//// 获取视频第一帧作为视频封面
		//imageFilePath, err := utils.GetSnapshot(savePath, name, 1)
		//if err != nil {
		//	publishVideosError(c, err.Error())
		//	return
		//}
		//if err := service.PostVideo(uid, savePath, imageFilePath, title); err != nil {
		//	publishVideosError(c, err.Error())
		//	return
		//}
		//publishVideosSuccess(c, "success")
	}
}

func publishVideosError(c *gin.Context, msg string) {
	c.JSON(http.StatusBadRequest, models.Response{
		StatusCode: -1,
		StatusMsg:  msg,
	})
}

func publishVideosSuccess(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, models.Response{
		StatusCode: 0,
		StatusMsg:  msg,
	})
}
