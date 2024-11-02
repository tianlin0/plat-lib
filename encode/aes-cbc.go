package encode

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/tianlin0/plat-lib/goroutines"
	"io"
)

func getAllKeyString(key string) string {
	keyLen := len(key)
	if keyLen >= 32 { //大于32的
		return key[0:32]
	}
	if keyLen >= 24 {
		return key[0:24]
	}
	if keyLen >= 16 {
		return key[0:16]
	}
	if keyLen > 0 && keyLen < 16 {
		for i := keyLen; i < 16; i++ {
			key += " " //不足的后面增加空格字符补齐
		}
		return key
	}
	return "jasonsjiang29121" //默认为GDP默认的key
}

// CBCDecrypt 解密
// decryptOri decrypt cypher text by key
func CBCDecrypt(cipherStr string, key string) (string, error) {
	var ciphertext []byte
	var err error

	key = getAllKeyString(key)
	keyByte := []byte(key)

	goroutines.GoSyncHandler(func(params ...interface{}) {
		var block cipher.Block
		ciphertext, err = hex.DecodeString(cipherStr)
		if err != nil {
			return
		}

		block, err = aes.NewCipher(keyByte)
		if err != nil {
			return
		}

		// The IV needs to be unique, but not secure. Therefore it's common to
		// include it at the beginning of the ciphertext.
		if len(ciphertext) < aes.BlockSize {
			err = fmt.Errorf("ciphertext too short")
			return
		}
		iv := ciphertext[:aes.BlockSize]
		ciphertext = ciphertext[aes.BlockSize:]

		stream := cipher.NewCBCDecrypter(block, iv)

		// XORKeyStream can work in-place if the two arguments are the same.
		stream.CryptBlocks(ciphertext, ciphertext)
	}, nil)

	if err != nil {
		return "", err
	}
	return string(bytes.TrimRight(ciphertext, string([]byte{0}))), nil
}

// CBCEncrypt 加密
// encryptOri encrypt plain text by key
func CBCEncrypt(plainStr string, key string) (string, error) {
	var ciphertext []byte
	var err error

	key = getAllKeyString(key)

	keyByte := []byte(key)

	goroutines.GoSyncHandler(func(params ...interface{}) {
		var block cipher.Block
		block, err = aes.NewCipher(keyByte)
		if err != nil {
			return
		}

		plaintext := []byte(plainStr)
		if len(plaintext)%aes.BlockSize != 0 {
			plaintext = append(plaintext, bytes.Repeat([]byte{0}, aes.BlockSize-len(plaintext)%aes.BlockSize)...)
		}
		// The IV needs to be unique, but not secure. Therefore it's common to
		// include it at the beginning of the ciphertext.
		ciphertext = make([]byte, aes.BlockSize+len(plaintext))
		iv := ciphertext[:aes.BlockSize]
		if _, err = io.ReadFull(rand.Reader, iv); err != nil {
			return
		}

		err = nil

		stream := cipher.NewCBCEncrypter(block, iv)
		stream.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)
	}, nil)

	if err != nil {
		return "", err
	}
	return hex.EncodeToString(ciphertext), nil
}
