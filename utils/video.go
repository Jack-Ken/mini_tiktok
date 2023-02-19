package utils

import (
	"fmt"
	"time"
)

func GetFileName(uid int64) string {
	// 生成唯一文件名称 时间戳+Uid
	t := time.Now().Unix()
	return fmt.Sprintf("%d-%d", t, uid)
}
