package encode

import (
	"bytes"
	"crypto"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
)

// RSA 默认
var RSA = &RSASecurity{}

// RSASecurity security
type RSASecurity struct {
	pubStr string          //公钥字符串
	priStr string          //私钥字符串
	pubkey *rsa.PublicKey  //公钥
	prikey *rsa.PrivateKey //私钥
}

// SetPublicKey 设置公钥
func (rsas *RSASecurity) SetPublicKey(pubStr string) (err error) {
	rsas.pubStr = pubStr
	rsas.pubkey, err = rsas.GetPublickey()
	return err
}

// SetPrivateKey 设置私钥
func (rsas *RSASecurity) SetPrivateKey(priStr string) (err error) {
	rsas.priStr = priStr
	rsas.prikey, err = rsas.GetPrivatekey()
	return err
}

// GetPrivatekey *rsa.PublicKey
func (rsas *RSASecurity) GetPrivatekey() (*rsa.PrivateKey, error) {
	return getPriKey([]byte(rsas.priStr))
}

// GetPublickey *rsa.PrivateKey
func (rsas *RSASecurity) GetPublickey() (*rsa.PublicKey, error) {
	return getPubKey([]byte(rsas.pubStr))
}

// PubKeyENCTYPT 公钥加密
func (rsas *RSASecurity) PubKeyENCTYPT(input []byte) ([]byte, error) {
	if rsas.pubkey == nil {
		return []byte(""), errors.New(`Please set the public key in advance`)
	}
	output := bytes.NewBuffer(nil)
	err := pubKeyIO(rsas.pubkey, bytes.NewReader(input), output, true)
	if err != nil {
		return []byte(""), err
	}
	return io.ReadAll(output)
}

// PubKeyDECRYPT 公钥解密
func (rsas *RSASecurity) PubKeyDECRYPT(input []byte) ([]byte, error) {
	if rsas.pubkey == nil {
		return []byte(""), errors.New(`Please set the public key in advance`)
	}
	output := bytes.NewBuffer(nil)
	err := pubKeyIO(rsas.pubkey, bytes.NewReader(input), output, false)
	if err != nil {
		return []byte(""), err
	}
	return io.ReadAll(output)
}

// PriKeyENCTYPT 私钥加密
func (rsas *RSASecurity) PriKeyENCTYPT(input []byte) ([]byte, error) {
	if rsas.prikey == nil {
		return []byte(""), errors.New(`Please set the private key in advance`)
	}
	output := bytes.NewBuffer(nil)
	err := priKeyIO(rsas.prikey, bytes.NewReader(input), output, true)
	if err != nil {
		return []byte(""), err
	}
	return io.ReadAll(output)
}

// PriKeyDECRYPT 私钥解密
func (rsas *RSASecurity) PriKeyDECRYPT(input []byte) ([]byte, error) {
	if rsas.prikey == nil {
		return []byte(""), errors.New(`Please set the private key in advance`)
	}
	output := bytes.NewBuffer(nil)
	err := priKeyIO(rsas.prikey, bytes.NewReader(input), output, false)
	if err != nil {
		return []byte(""), err
	}

	return io.ReadAll(output)
}

// SignMd5WithRsa 使用RSAWithMD5算法签名
func (rsas *RSASecurity) SignMd5WithRsa(data string) (string, error) {
	md5Hash := md5.New()
	s_data := []byte(data)
	md5Hash.Write(s_data)
	hashed := md5Hash.Sum(nil)

	signByte, err := rsa.SignPKCS1v15(rand.Reader, rsas.prikey, crypto.MD5, hashed)
	sign := base64.StdEncoding.EncodeToString(signByte)
	return sign, err
}

// SignSha1WithRsa 使用RSAWithSHA1算法签名
func (rsas *RSASecurity) SignSha1WithRsa(data string) (string, error) {
	sha1Hash := sha1.New()
	s_data := []byte(data)
	sha1Hash.Write(s_data)
	hashed := sha1Hash.Sum(nil)

	signByte, err := rsa.SignPKCS1v15(rand.Reader, rsas.prikey, crypto.SHA1, hashed)
	sign := base64.StdEncoding.EncodeToString(signByte)
	return sign, err
}

// SignSha256WithRsa 使用RSAWithSHA256算法签名
func (rsas *RSASecurity) SignSha256WithRsa(data string) (string, error) {
	sha256Hash := sha256.New()
	s_data := []byte(data)
	sha256Hash.Write(s_data)
	hashed := sha256Hash.Sum(nil)

	signByte, err := rsa.SignPKCS1v15(rand.Reader, rsas.prikey, crypto.SHA256, hashed)
	sign := base64.StdEncoding.EncodeToString(signByte)
	return sign, err
}

// VerifySignMd5WithRsa 使用RSAWithMD5验证签名
func (rsas *RSASecurity) VerifySignMd5WithRsa(data string, signData string) error {
	sign, err := base64.StdEncoding.DecodeString(signData)
	if err != nil {
		return err
	}
	hash := md5.New()
	hash.Write([]byte(data))
	return rsa.VerifyPKCS1v15(rsas.pubkey, crypto.MD5, hash.Sum(nil), sign)
}

// VerifySignSha1WithRsa 使用RSAWithSHA1验证签名
func (rsas *RSASecurity) VerifySignSha1WithRsa(data string, signData string) error {
	sign, err := base64.StdEncoding.DecodeString(signData)
	if err != nil {
		return err
	}
	hash := sha1.New()
	hash.Write([]byte(data))
	return rsa.VerifyPKCS1v15(rsas.pubkey, crypto.SHA1, hash.Sum(nil), sign)
}

// VerifySignSha256WithRsa 使用RSAWithSHA256验证签名
func (rsas *RSASecurity) VerifySignSha256WithRsa(data string, signData string) error {
	sign, err := base64.StdEncoding.DecodeString(signData)
	if err != nil {
		return err
	}
	hash := sha256.New()
	hash.Write([]byte(data))

	return rsa.VerifyPKCS1v15(rsas.pubkey, crypto.SHA256, hash.Sum(nil), sign)
}
