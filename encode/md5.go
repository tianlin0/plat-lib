package encode

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
)

// Md5 计算出md5的值
func Md5(s string) string {
	d := md5.Sum([]byte(s))
	return hex.EncodeToString(d[:])
}

// Sha256Hex 转换为sha256字符
func Sha256Hex(s string) string {
	bytes2 := sha256.Sum256([]byte(s))
	return hex.EncodeToString(bytes2[:])
}

// HmacSha256 转化为hmac
func HmacSha256(s, secret string) string {
	hashed := hmac.New(sha256.New, []byte(secret))
	hashed.Write([]byte(s))
	return base64.StdEncoding.EncodeToString(hashed.Sum(nil))
}
