package encode

import (
	"encoding/base64"
)

// XorEncode 异或加密
func XorEncode(msg string, key int) string {
	byteList := []byte(msg)
	pwd := make([]byte, len(byteList))
	for i := 0; i < len(byteList); i++ {
		pwd[i] = byteList[i] ^ byte(key)
	}
	return base64.StdEncoding.EncodeToString(pwd)
}

// XorDecode 异或解密
func XorDecode(msg string, key int) string {
	pwdList, _ := base64.StdEncoding.DecodeString(msg)
	old := make([]byte, len(pwdList))
	for i := 0; i < len(pwdList); i++ {
		old[i] = pwdList[i] ^ byte(key)
	}
	return string(old)
}
