package utils

import "C"

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"go.uber.org/zap"

	"github.com/disintegration/imaging"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

/*
ffmpeg获取视频第一帧的图片
*/

func GetSnapshot(videoPath string, imageName string, frameNum int) (ImagePath string, err error) {
	snapshotPath := filepath.Join("../static/images", imageName)
	buf := bytes.NewBuffer(nil)
	err = ffmpeg.Input(videoPath).Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", frameNum)}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(buf, os.Stdout).
		Run()

	if err != nil {
		zap.L().Error("生成缩略图失败1：", zap.Error(err))
		return "", err
	}

	img, err := imaging.Decode(buf)
	if err != nil {
		zap.L().Error("生成缩略图失败2：", zap.Error(err))
		return "", err
	}

	err = imaging.Save(img, snapshotPath+".png")
	if err != nil {
		zap.L().Error("生成缩略图失败3：", zap.Error(err))
		return "", err
	}

	imgPath := snapshotPath + ".png"
	return imgPath, nil
}
