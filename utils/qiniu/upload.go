package qiniu

import (
	"bytes"
	"context"
	"net/http"
	"strings"

	config "mini_tiktok/internal/initialize"

	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

// 上传图片到七牛云，然后返回状态和图片的url
func UploadimageToQiNiu(data []byte, name string) (string, error) {
	cnf := config.Conf.QiNiuConfig
	var AccessKey = cnf.Accesskey // 秘钥对
	var SerectKey = cnf.Sercetkey
	var Bucket = cnf.Bucket      // 空间名称
	var ImgUrl = cnf.Qiniuserver // 自定义域名或测试域名

	// 根据userid创建视频名称
	//name := utils.GetFileName(uid)
	// 检查文件类型 MIME 的功能，例如image/jpeg , video/mp4
	contentType := http.DetectContentType(data)

	t := strings.Split(contentType, "/")
	root, suffix := t[0], t[1]

	src := bytes.NewReader(data)

	putPlicy := storage.PutPolicy{
		Scope: Bucket,
	}
	mac := qbox.NewMac(AccessKey, SerectKey)

	// 获取上传凭证
	upToken := putPlicy.UploadToken(mac)

	// 配置参数
	cfg := storage.Config{
		Zone:          &storage.ZoneHuabei, // 华北区
		UseCdnDomains: false,
		UseHTTPS:      false, // 非https
	}
	formUploader := storage.NewFormUploader(&cfg)

	ret := storage.PutRet{}        // 上传后返回的结果
	putExtra := storage.PutExtra{} // 额外参数

	// 上传 自定义key，可以指定上传目录及文件名和后缀，
	key := root + "/" + name + "." + suffix // 上传路径，如果当前目录中已存在相同文件，则返回上传失败错误
	err := formUploader.Put(context.Background(), &ret, upToken, key, src, int64(len(data)), &putExtra)

	if err != nil {
		return "", err
	}

	url := ImgUrl + "/" + ret.Key // 返回上传后的文件访问路径
	return url, nil
}
