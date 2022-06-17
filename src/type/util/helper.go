package util

import (
	"crypto/md5"
	"encoding/hex"
)

func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	hashStr := hex.EncodeToString(hash[:])
	return hashStr
}
