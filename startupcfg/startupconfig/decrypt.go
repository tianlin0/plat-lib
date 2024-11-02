package startupconfig

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"

	"gopkg.in/yaml.v3"
)

var key = []byte{12, 7, 21, 9, 8, 21, 12, 15, 7, 8, 1, 84, 95, 87, 84, 87}

func init() {
	for i := range key {
		key[i] = key[i] ^ 0x66
	}
}

// DecDecrypt 解密
func DecDecrypt(cipherStr string) (string, error) {
	return decryptOri(cipherStr, key)
}

// decryptOri decrypt cypher text by key
func decryptOri(cipherStr string, key []byte) (string, error) {
	defer func() {
		if err := recover(); err != nil {
			log.Fatalln("panic: ", err)
		}
	}()

	ciphertext, err := hex.DecodeString(cipherStr)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(ciphertext) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCBCDecrypter(block, iv)

	// XORKeyStream can work in-place if the two arguments are the same.
	stream.CryptBlocks(ciphertext, ciphertext)
	return string(bytes.TrimRight(ciphertext, string([]byte{0}))), nil
}

// Encrypt 加密
func Encrypt(plainStr string) (string, error) {
	return encryptOri(plainStr, key)
}

// encryptOri encrypt plain text by key
func encryptOri(plainStr string, key []byte) (string, error) {
	defer func() {
		if err := recover(); err != nil {
			log.Fatalln("panic: ", err)
		}
	}()

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	plaintext := []byte(plainStr)
	if len(plaintext)%aes.BlockSize != 0 {
		plaintext = append(plaintext, bytes.Repeat([]byte{0}, aes.BlockSize-len(plaintext)%aes.BlockSize)...)
	}
	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCBCEncrypter(block, iv)
	stream.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)

	// It's important to remember that ciphertexts must be authenticated
	// (i.e. by using crypto/hmac) as well as being encrypted in order to
	// be secure.
	return hex.EncodeToString(ciphertext), nil
}

// Encrypted 加密串儿
type Encrypted string

// Get 获取解密串儿
func (e Encrypted) Get() (string, error) {
	if e == "" {
		return "", nil
	}
	str, err := DecDecrypt(string(e))
	if err != nil {
		return "", err
	}
	return str, nil
}

// MarshalJSON 实现json Marshaler接口 自定义json 编码
func (e *Encrypted) MarshalJSON() ([]byte, error) {
	decrypted, err := e.Get()
	if err != nil {
		return nil, err
	}
	return ([]byte)(fmt.Sprintf("\"%s\"", decrypted)), nil
}

// MarshalYAML 实现yaml Marshaler接口
func (e *Encrypted) MarshalYAML() (interface{}, error) {
	decrypted, err := e.Get()
	if err != nil {
		return nil, err
	}
	return decrypted, nil
}

// Decrypted 解密串儿
type Decrypted string

// UnmarshalJSON 实现Unmarshaler接口 自定义json解码
func (d *Decrypted) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return fmt.Errorf("Unmarshal Decrypted(%s) to string failed: %w ", data, err)
	}
	if str == "" {
		return nil
	}
	dd, err := DecDecrypt(str)
	if err != nil {
		return fmt.Errorf("DecDecrypt Decrypted(%s) failed: %w ", str, err)
	}
	*d = Decrypted(dd)
	return nil
}

// UnmarshalYAML 实现Unmarshaler接口 自定义yaml解码
func (d *Decrypted) UnmarshalYAML(value *yaml.Node) error {
	if value.Value == "" {
		return nil
	}
	dd, err := DecDecrypt(value.Value)
	if err != nil {
		return fmt.Errorf("DecDecrypt Decrypted(%s) failed: %w ", value.Value, err)
	}
	*d = Decrypted(dd)
	return nil
}

// String
func (d Decrypted) String() string {
	return string(d)
}
