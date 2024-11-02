package encode

import (
	"encoding/base64"
	ecb "github.com/haowanxing/go-aes-ecb"
)

//电码本模式（Electronic Codebook Book (ECB)）
type ecbParam struct {
	pad   string
	block int
}

func encode_7_128Block(origin, key string) (string, error) {
	ciphertext := ecb.PKCS7Padding([]byte(origin), 16)

	crypted, err := ecb.AesEncrypt(ciphertext, []byte(key))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(crypted), nil
}

func decode_7_128Block(crypted, key string) (string, error) {
	cryptedByte, err := base64.StdEncoding.DecodeString(crypted)
	if err != nil {
		return "", err
	}
	origin, err := ecb.AesDecrypt(cryptedByte, []byte(key))
	if err != nil {
		return "", err
	}
	origin = ecb.PKCS7UnPadding(origin)
	return string(origin), nil
}
func encode_0_128Block(origin, key string) (string, error) {
	ciphertext := ecb.ZerosPadding([]byte(origin), 16)

	crypted, err := ecb.AesEncrypt(ciphertext, []byte(key))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(crypted), nil
}

func decode_0_128Block(crypted, key string) (string, error) {
	cryptedByte, err := base64.StdEncoding.DecodeString(crypted)
	if err != nil {
		return "", err
	}
	origin, err := ecb.AesDecrypt(cryptedByte, []byte(key))
	if err != nil {
		return "", err
	}
	origin = ecb.ZerosUnPadding(origin)
	return string(origin), nil
}
func encode_7_192Block(origin, key string) (string, error) {
	ciphertext := ecb.PKCS7Padding([]byte(origin), 16)

	crypted, err := ecb.AesEncrypt(ciphertext, []byte(key))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(crypted), nil
}

func decode_7_192Block(crypted, key string) (string, error) {
	cryptedByte, err := base64.StdEncoding.DecodeString(crypted)
	if err != nil {
		return "", err
	}
	origin, err := ecb.AesDecrypt(cryptedByte, []byte(key))
	if err != nil {
		return "", err
	}
	origin = ecb.PKCS7UnPadding(origin)
	return string(origin), nil
}
