package encode

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
)

// PublicEncrypt 公钥加密
func PublicEncrypt(data, publicKey string) (string, error) {
	grsa := rsaSecurity{}
	_ = grsa.SetPublicKey(publicKey)

	rsadata, err := grsa.PubKeyENCTYPT([]byte(data))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(rsadata), nil
}

// PriKeyEncrypt 私钥加密
func PriKeyEncrypt(data, privateKey string) (string, error) {
	grsa := rsaSecurity{}
	_ = grsa.SetPrivateKey(privateKey)
	rsadata, err := grsa.PriKeyENCTYPT([]byte(data))
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(rsadata), nil
}

// PublicDecrypt 公钥解密
func PublicDecrypt(data, publicKey string) (string, error) {
	databs, _ := base64.StdEncoding.DecodeString(data)
	grsa := rsaSecurity{}
	_ = grsa.SetPublicKey(publicKey)

	rsadata, err := grsa.PubKeyDECRYPT([]byte(databs))
	if err != nil {
		return "", err
	}
	return string(rsadata), nil

}

// PriKeyDecrypt 私钥解密
func PriKeyDecrypt(data, privateKey string) (string, error) {
	databs, _ := base64.StdEncoding.DecodeString(data)

	grsa := rsaSecurity{}
	_ = grsa.SetPrivateKey(privateKey)

	rsadata, err := grsa.PriKeyDECRYPT([]byte(databs))
	if err != nil {
		return "", err
	}
	return string(rsadata), nil
}

// RSAPublicEncrypt 公钥加密
func RSAPublicEncrypt(encryptStr string, publicKeyStr string) (string, error) {
	// -----BEGIN PUBLIC KEY-----
	buf := []byte(publicKeyStr)

	// pem 解码
	block, _ := pem.Decode(buf)

	// x509 解码
	publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", err
	}

	// 类型断言
	publicKey := publicKeyInterface.(*rsa.PublicKey)

	//对明文进行加密
	encryptedStr, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, []byte(encryptStr))
	if err != nil {
		return "", err
	}

	//返回密文
	return base64.URLEncoding.EncodeToString(encryptedStr), nil
}

// RSAPrivateDecrypt 私钥解密
func RSAPrivateDecrypt(decryptStr string, privateKeyStr string) (string, error) {
	//-----BEGIN RSA PRIVATE KEY-----
	buf := []byte(privateKeyStr)

	// pem 解码
	block, _ := pem.Decode(buf)

	// X509 解码
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}
	decryptBytes, err := base64.URLEncoding.DecodeString(decryptStr)

	//对密文进行解密
	decrypted, _ := rsa.DecryptPKCS1v15(rand.Reader, privateKey, decryptBytes)

	//返回明文
	return string(decrypted), nil
}
