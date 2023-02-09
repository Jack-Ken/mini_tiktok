package utils

import (
	"crypto/md5"
	"encoding/hex"
)

var secret = "pbkdf2-sha512$%s$%s"

func EncryptPassword(oPassword string) string {
	h := md5.New()
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum([]byte(oPassword)))
}
