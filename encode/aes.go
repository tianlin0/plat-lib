package encode

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"github.com/tianlin0/plat-lib/conv"
	"strings"
)

func pKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padText...)
}

func pKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unPadding := int(origData[length-1])
	if length-unPadding < 0 { // 不能<0
		return []byte{}
	}
	return origData[:(length - unPadding)]
}

func aesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	origData = pKCS5Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypt := make([]byte, len(origData))
	blockMode.CryptBlocks(crypt, origData)
	return crypt, nil
}

func aesDecrypt(crypt, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(crypt))
	blockMode.CryptBlocks(origData, crypt)
	origData = pKCS5UnPadding(origData)
	return origData, nil
}

// AesEncryptBase64 对字符串对称加密，输入输出都为base64字符
func AesEncryptBase64(origString string, keyBase64 string) (string, error) {
	if origString == "" {
		return "", nil
	}
	key, err := base64.StdEncoding.DecodeString(keyBase64)
	if err != nil {
		//表示传的不是base64格式的，则内部兼容处理一下
		key = []byte(keyBase64)
	}

	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		//太短或太长，则进行转化
		newKey := Md5(string(key))
		key = []byte(newKey)
	}

	origData := []byte(origString)

	encode, err := aesEncrypt(origData, key)
	if err != nil {
		return "", err
	}

	retBase64 := base64.StdEncoding.EncodeToString(encode)
	return retBase64, nil
}

// AesDecryptBase64 对字符串对称解密，输入输出都为base64字符
func AesDecryptBase64(encodeBase64 string, keyBase64 string) (string, error) {
	if encodeBase64 == "" {
		return "", nil
	}

	key, err := base64.StdEncoding.DecodeString(keyBase64)
	if err != nil {
		key = []byte(keyBase64)
	}

	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		//太短或太长，则进行转化
		newKey := Md5(string(key))
		key = []byte(newKey)
	}

	encodeByte, err := base64.StdEncoding.DecodeString(encodeBase64)
	if err != nil {
		return "", err
	}

	decodeByte, err := aesDecrypt(encodeByte, key)
	if err != nil {
		return "", err
	}

	return string(decodeByte), nil
}

// Encrypt aes
func Encrypt(encryptStr string, key []byte, iv string) (string, error) {
	encryptBytes := []byte(encryptStr)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	blockSize := block.BlockSize()
	encryptBytes = pKCS5Padding(encryptBytes, blockSize)

	blockMode := cipher.NewCBCEncrypter(block, []byte(iv))
	encrypted := make([]byte, len(encryptBytes))
	blockMode.CryptBlocks(encrypted, encryptBytes)
	return base64.URLEncoding.EncodeToString(encrypted), nil
}

// Decrypt aes
func Decrypt(decryptStr string, key []byte, iv string) (string, error) {
	decryptBytes, err := base64.URLEncoding.DecodeString(decryptStr)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	blockMode := cipher.NewCBCDecrypter(block, []byte(iv))
	decrypted := make([]byte, len(decryptBytes))

	blockMode.CryptBlocks(decrypted, decryptBytes)
	decrypted = pKCS5UnPadding(decrypted)
	return string(decrypted), nil
}

// AesDecryptString 解密
func AesDecryptString(encode string, key string) string {
	encode = strings.ReplaceAll(encode, "\r", "")
	encode = strings.ReplaceAll(encode, "\n", "")
	encode = strings.ReplaceAll(encode, "\t", "")
	encode = strings.ReplaceAll(encode, " ", "")

	newString, err := AesDecryptBase64(encode, key)
	if newString != "" && err == nil {
		return newString
	}
	return ""
}

// AesEncryptString 加密
func AesEncryptString(obj interface{}, key string) string {
	jsStr := conv.String(obj)
	encode, err := AesEncryptBase64(jsStr, key)
	if encode != "" && err == nil {
		return encode
	}
	return ""
}
