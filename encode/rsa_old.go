package encode

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"strings"
)

/**
  pubKey, priKey, _ := utils.GenRsaPublicAndPrivateKey()

  fmt.Println(pubKey)
  fmt.Println(priKey)

  pubSecret, _ := utils.RsaEncryptByPublicKey("我是谁", pubKey)
  fmt.Println(pubSecret)

  pubOld, _ := utils.RsaDecryptByPrivateKey(pubSecret, priKey)
  fmt.Println(pubOld)

  priSecret, _ := utils.RsaEncryptByPrivateKey("my name is", priKey)
  fmt.Println(priSecret)

  priOld, _ := utils.RsaDecryptByPublicKey(priSecret, pubKey)
  fmt.Println(priOld)
*/

// GenRsaPublicAndPrivateKey111 RSA公钥私钥产生
func GenRsaPublicAndPrivateKey111() (pubKey, priKey string, err error) {
	bits := 1024
	// 生成私钥文件
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return "", "", err
	}
	return encodePKCS1(privateKey)
}

type RsaLength int
type RsaFormat string

// GenRsaPublicAndPrivateKey 老的方法
func GenRsaPublicAndPrivateKey(rsaLength RsaLength, rsaFormat RsaFormat, rsaPass string) (
	pubKey, priKey string, err error) {
	var bitLength RsaLength = 2048

	rsaLengthList := []RsaLength{512, 1024, 2048, 4096}
	if rsaLength > 0 {
		for _, one := range rsaLengthList {
			if one == rsaLength {
				bitLength = one
				break
			}
		}
	}

	// 生成私钥文件
	privateKey, err := rsa.GenerateKey(rand.Reader, int(bitLength))
	if err != nil {
		return "", "", err
	}

	if rsaFormat == "PKCS#8" {
		return encodePKCS8(privateKey)
	} else if rsaFormat == "PKCS#1" {
		return encodePKCS1(privateKey)
	}
	return "", "", fmt.Errorf("rsaFormat error")
}

func encodePKCS1(privateKey *rsa.PrivateKey) (string, string, error) {
	derStream := x509.MarshalPKCS1PrivateKey(privateKey)
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derStream,
	}

	var val bytes.Buffer
	var key bytes.Buffer
	err := pem.Encode(&val, block)
	if err != nil {
		return "", "", err
	}
	// 生成公钥文件
	publicKey := &privateKey.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return "", "", err
	}
	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPkix,
	}
	err = pem.Encode(&key, block)
	if err != nil {
		return "", "", err
	}

	return key.String(), val.String(), nil
}

func encodePKCS8(privateKey *rsa.PrivateKey) (string, string, error) {
	derStream, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return "", "", err
	}

	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derStream,
	}

	var val bytes.Buffer
	var key bytes.Buffer
	err = pem.Encode(&val, block)
	if err != nil {
		return "", "", err
	}
	// 生成公钥文件
	publicKey := &privateKey.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return "", "", err
	}
	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPkix,
	}
	err = pem.Encode(&key, block)
	if err != nil {
		return "", "", err
	}

	return key.String(), val.String(), nil
}

var (
	publickeyStart = "-----BEGIN PUBLIC KEY-----"
	publickeyEnd   = "-----END PUBLIC KEY-----"
	privateStart   = "-----BEGIN RSA PRIVATE KEY-----"
	privateEnd     = "-----END RSA PRIVATE KEY-----"
)

func buitPubKey(pubKey string) string {
	if !strings.HasPrefix(pubKey, publickeyStart) {
		pubKey = publickeyStart + "\n" + pubKey
	}
	if !strings.HasSuffix(pubKey, publickeyEnd) {
		pubKey = pubKey + "\n" + publickeyEnd
	}
	return pubKey
}
func buitPriKey(priKey string) string {
	if !strings.HasPrefix(priKey, privateStart) {
		priKey = privateStart + "\n" + priKey
	}
	if !strings.HasSuffix(priKey, privateEnd) {
		priKey = priKey + "\n" + privateEnd
	}
	return priKey
}

// RsaEncryptByPublicKey 公钥加密
func RsaEncryptByPublicKey(oldStr string, pubKey string) (string, error) {
	pubKey = buitPubKey(pubKey)
	rsaClient := &RSASecurity{}
	rsaClient.SetPublicKey(pubKey)
	pubenctypt, err := rsaClient.PubKeyENCTYPT([]byte(oldStr))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(pubenctypt), nil
}

// RSA_Encrypt 非对称加密
func RSA_Encrypt(plainText string, pubKey string) string {
	pubKey = buitPubKey(pubKey)

	//pem解码
	block, _ := pem.Decode([]byte(pubKey))
	//x509解码

	publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		panic(interface{}(err))
	}
	//类型断言
	publicKey := publicKeyInterface.(*rsa.PublicKey)
	//对明文进行加密
	cipherText, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, []byte(plainText))
	if err != nil {
		panic(interface{}(err))
	}
	//返回密文
	return base64.StdEncoding.EncodeToString(cipherText)
}

// RSA_Decrypt 非对称解密
func RSA_Decrypt(cipherText string, priKey string) []byte {
	priKey = buitPriKey(priKey)

	//pem解码
	block, _ := pem.Decode([]byte(priKey))
	//X509解码
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		panic(interface{}(err))
	}

	cipherTextByte, err := base64.StdEncoding.DecodeString(cipherText)
	//对密文进行解密
	plainText, _ := rsa.DecryptPKCS1v15(rand.Reader, privateKey, cipherTextByte)
	//返回明文
	return plainText
}

// RsaDecryptByPrivateKey 私钥解密
func RsaDecryptByPrivateKey(encodeStr string, priKey string) (string, error) {
	priKey = buitPriKey(priKey)
	encodeByte, err := base64.StdEncoding.DecodeString(encodeStr)
	if err != nil {
		encodeByte = []byte(encodeStr)
	}
	rsaClient := &RSASecurity{}
	_ = rsaClient.SetPrivateKey(priKey)
	pridecrypt, err := rsaClient.PriKeyDECRYPT(encodeByte)
	if err != nil {
		return "", err
	}
	return string(pridecrypt), nil
}

// RsaDecryptByPublicKey 公钥解密
func RsaDecryptByPublicKey(encodeStr string, pubKey string) (string, error) {
	pubKey = buitPubKey(pubKey)
	encodeByte, err := base64.StdEncoding.DecodeString(encodeStr)
	if err != nil {
		encodeByte = []byte(encodeStr)
	}
	rsaClient := &RSASecurity{}
	_ = rsaClient.SetPublicKey(pubKey)
	pubdecrypt, err := rsaClient.PubKeyDECRYPT(encodeByte)
	if err != nil {
		return "", err
	}
	return string(pubdecrypt), nil
}

// RsaEncryptByPrivateKey 私钥加密
func RsaEncryptByPrivateKey(oldStr string, priKey string) (string, error) {
	priKey = buitPriKey(priKey)
	rsaClient := &RSASecurity{}
	rsaClient.SetPrivateKey(priKey)
	printCtypt, err := rsaClient.PriKeyENCTYPT([]byte(oldStr))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(printCtypt), nil
}
