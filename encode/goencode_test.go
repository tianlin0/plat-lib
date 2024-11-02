package encode

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"os"

	"testing"
)

var Pubkey = `-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDhHrkiU2FB137Yl9LbwFKcRxqg
UywnhAkWW28KH1impnt7Ir1a71JhR1NFaUlnk+rG9X6cpDy4bS9I4g07/f8ULLsv
cNVDtcTptYNAkczQ8Gozvr+tkNmI9D+s0tFa0vVoq2n10NLtg7YdF2WM0M29xlEE
cQ+Kk42eEEBl4Rk5YQIDAQAB
-----END PUBLIC KEY-----`

var Pirvatekey = `-----BEGIN PRIVATE KEY-----
MIICdgIBADANBgkqhkiG9w0BAQEFAASCAmAwggJcAgEAAoGBAOEeuSJTYUHXftiX
0tvAUpxHGqBTLCeECRZbbwofWKame3sivVrvUmFHU0VpSWeT6sb1fpykPLhtL0ji
DTv9/xQsuy9w1UO1xOm1g0CRzNDwajO+v62Q2Yj0P6zS0VrS9WirafXQ0u2Dth0X
ZYzQzb3GUQRxD4qTjZ4QQGXhGTlhAgMBAAECgYBt5ImLcBhyA7gwEy0jiObK0wr0
aKWNRK8K8udpkZO9BlgQ7Axzb5BPXHoR0Cu9HD/nj7+Wx7W8cdA7S94aAwuY6stH
YRlVkeAUJBsl0ur+bu15cMFL4NA7+bg7FQXUzElp7LyefPcMU1FcPEj19QxOKv4X
cKVwfm8fWlzQzSM/YQJBAP8+Q/7rKbZGFZk0nI4XiiPlOr4P//8EclO1nBCHdCnL
AApxK6E9cDTbTPRHTWbLQMToMbFuqe3Rrhi7OrqKXWMCQQDhyZfu3jSJ+xl4bpG4
rMBp96+38TF9YNxrwagtw6UXamAAXgmMjg2lKrVlt5yzi8bitjL5hhwWV2LMuWFb
RltrAkEA8hVOTGMiRrymE47wxVvSK0Vot4dZV7gR7w8anBq8tD7TJRQ9O0qYN6mf
jThrUwmHvrozF4RMK0FqDA7YHsDI3QJAVPY78smwsX9IbVYGBZ0T5owqlifvfIN3
TiEYPOhS9kW0DE9WfopxvgYdLkJyd+mQFH2FHvoFFa8aYXkclnEaMwJAXipJkuQX
GWLLZF6GjTGMTSy2i5QE3ZXIidDH4Y67dcGVdCm+SgF55VYa4q8emGt9GgMuLQWx
2V4NAJXYx7HcNw==
-----END PRIVATE KEY-----`

// 公钥解密私钥加密
//func TestApplyPriEPub(t *testing.T) {
//	str := "hello world"
//	prienctypt, err := gorsa.PriKeyEncrypt(str, Pirvatekey)
//	if err != nil {
//		fmt.Println("err1", err)
//		return
//	}
//
//	fmt.Println("============")
//	fmt.Println(prienctypt)
//	fmt.Println("============")
//
//	pubdecrypt, err := gorsa.PublicDecrypt(prienctypt, Pubkey)
//	if err != nil {
//		fmt.Println("err2", err)
//		return
//	}
//	if string(pubdecrypt) != str {
//		fmt.Println("error")
//		return
//	}
//	fmt.Println("success")
//	return
//}

func TestXorEncode(t *testing.T) {
	pubenctypt, _ := PublicEncrypt(`hello world`, Pubkey)
	fmt.Println(pubenctypt)

	pubenctypt, _ = PublicEncrypt(`hello world`, Pubkey)
	fmt.Println(pubenctypt)

	pubenctypt, _ = RsaEncryptByPublicKey(`hello world`, Pubkey)
	fmt.Println(pubenctypt)
}
func TestAes(t *testing.T) {
	fmt.Println(AesEncryptBase64("欢迎来到chacuo.net", "1234"))
}
func TestEcb(t *testing.T) {
	content := "欢迎来到chacuo"
	key := "0123456789abcdef11111112"

	aa, _ := encode_7_192Block(content, key)

	fmt.Println(aa)

	fmt.Println(decode_7_192Block(aa, key))

	aa, _ = encode_0_128Block(content, key)
	fmt.Println(aa)

	fmt.Println(decode_0_128Block(aa, key))
}

func TestEcb2(t *testing.T) {
	//key := "0123456789abcdef"
	////加密
	//ciphertext := security.Encrypt(security.PKCS7Pad([]byte("欢迎来到chacuo")), key)
	////解密
	//plaintext := string(security.PKCS7UPad([]byte(security.Decrypt(ciphertext, key))))
	//
	//fmt.Println(base64.StdEncoding.EncodeToString(ciphertext), plaintext)
}

func TestAAA(t *testing.T) {
	priKey := `-----BEGIN PRIVATE KEY-----
MIICdgIBADANBgkqhkiG9w0BAQEFAASCAmAwggJcAgEAAoGBALSfh8nVLtXIhSTN
c8Mkk4m1KaxRvqg7GZ0wg368kqiI676SG5cgmjZdjAAuK4pBoa0VJcee2gtqo14C
bc4iEUfy0IfDMCYlBzQtvGQ9saf50NYB60eHLA1dbLvUc1RHVdHPUXzpE3s0RLWM
5nahQwL/fRt7oxgDcoZ8qRdtEllBAgMBAAECgYBDn1ZfIgke0KvIU4L7lD4IWGL5
uMEAit/UEc2pLUBbCKf5+QmLUxFpOSypBKAYaun0uu4iBj7r90iicZZajjaZcCOv
f26CBZqkm4qKwt/AhYYcart4mNuW+5AnvM0xGAKnYJc8FzHyAomnCFAoj7ei/OBF
45zE2OzE1GyrWn2IAQJBAN0UX4NDr0t0fV3gRQyyXFNEMp9NTExcwQTMvO2FqLc7
/7lscE4baDyt9cW0dgGzytoM7sKi4Mn0S5tXMagkUhECQQDRJ0CBt90OhwFebyD4
oRtZsEJinkfPLJ5SX7CjxIq0iaz9yV9BmBJrZ+gSvk4Vf9UwK/OW1Y7kK7PQhucd
EGQxAkEAuWEd/gnBcboKba9i9xSQilnDQQUmF1onmAi920WahZtQAYHGYhhlPYx5
bAC4exDx5gm2I4tEhtPMmkNxJhbeoQJAVglOiM3omkRA9ObD6mLjjFZsSIMRyRBy
pDIGyKdd43xK9C71B1eWJCafGa69Eiz+to0t69s3p3auxlXoFlWa0QJARljTi4X/
NrhtBDc4bcq25SUBi4Uv3uZDQTSuJp1nCyUliJTMw0Pasfo7msQIfkv3JH87Ws/x
gTjci7afvT49rw==
-----END PRIVATE KEY-----
`
	aaa, err := EncryptPEM([]byte(priKey), []byte("1"))
	fmt.Println(base64.StdEncoding.EncodeToString(aaa.Bytes), err)

	return

	// Generate RSA Keys
	miryanPrivateKey, err := rsa.GenerateKey(rand.Reader, 2048)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	miryanPublicKey := &miryanPrivateKey.PublicKey

	raulPrivateKey, err := rsa.GenerateKey(rand.Reader, 2048)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	raulPublicKey := &raulPrivateKey.PublicKey

	fmt.Println("Private Key : ", miryanPrivateKey)
	fmt.Println("Public key ", miryanPublicKey)
	fmt.Println("Private Key : ", raulPrivateKey)
	fmt.Println("Public key ", raulPublicKey)

	//Encrypt Miryan Message
	message := []byte("the code must be like a piece of music")
	label := []byte("")
	hash := sha256.New()

	ciphertext, err := rsa.EncryptOAEP(hash, rand.Reader, raulPublicKey, message, label)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("OAEP encrypted [%s] to \n[%x]\n", string(message), ciphertext)
	fmt.Println()

	// Message - Signature
	var opts rsa.PSSOptions
	opts.SaltLength = rsa.PSSSaltLengthAuto // for simple example
	PSSmessage := message
	newhash := crypto.SHA256
	pssh := newhash.New()
	pssh.Write(PSSmessage)
	hashed := pssh.Sum(nil)

	signature, err := rsa.SignPSS(rand.Reader, miryanPrivateKey, newhash, hashed, &opts)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("PSS Signature : %x\n", signature)

	// Decrypt Message
	plainText, err := rsa.DecryptOAEP(hash, rand.Reader, raulPrivateKey, ciphertext, label)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("OAEP decrypted [%x] to \n[%s]\n", ciphertext, plainText)

	//Verify Signature
	err = rsa.VerifyPSS(miryanPublicKey, newhash, hashed, signature, &opts)

	if err != nil {
		fmt.Println("Who are U? Verify Signature failed")
		os.Exit(1)
	} else {
		fmt.Println("Verify Signature successful")
	}
}

func PrivateKeyToEncryptedPEM(privateKey string, pwd string) ([]byte, error) {

	block, _ := pem.Decode([]byte(privateKey))

	// Encrypt the pem

	if pwd != "" {

		blocks, err := x509.EncryptPEMBlock(rand.Reader, block.Type, block.Bytes, []byte(pwd), x509.PEMCipherAES256)

		if err != nil {

			return nil, err

		}

		var pemOut bytes.Buffer
		pem.Encode(&pemOut, blocks)

		fmt.Println(pemOut.String())

		block = blocks

	}

	return pem.EncodeToMemory(block), nil

}

func DecryptPEM(pemRaw []byte, passwd []byte) ([]byte, error) {
	block, _ := pem.Decode(pemRaw)
	if block == nil {
		return nil, fmt.Errorf("Failed decoding PEM. Block must be different from nil. [% x]", pemRaw)
	}

	if !x509.IsEncryptedPEMBlock(block) {
		return nil, fmt.Errorf("Failed decryptPEM PEM. it's not a decryped PEM [%s]", pemRaw)
	}

	der, err := x509.DecryptPEMBlock(block, passwd)
	if err != nil {
		return nil, fmt.Errorf("Failed PEM decryption [%s]", err)
	}

	privateKey, err := DERToPrivateKey(der)
	if err != nil {
		return nil, err
	}

	var raw []byte
	switch k := privateKey.(type) {
	case *ecdsa.PrivateKey:
		raw, err = x509.MarshalECPrivateKey(k)
		if err != nil {
			return nil, err
		}
	case *rsa.PrivateKey:
		raw = x509.MarshalPKCS1PrivateKey(k)
	default:
		return nil, fmt.Errorf("Invalid key type. It must be *ecdsa.PrivateKey or *rsa.PrivateKey")
	}

	rawBase64 := base64.StdEncoding.EncodeToString(raw)
	derBase64 := base64.StdEncoding.EncodeToString(der)
	if rawBase64 != derBase64 {
		return nil, fmt.Errorf("Invalid decrypt PEM: raw does not match with der")
	}

	block = &pem.Block{
		Type:  block.Type,
		Bytes: der,
	}

	return pem.EncodeToMemory(block), nil
}

func DERToPrivateKey(der []byte) (key interface{}, err error) {
	if key, err = x509.ParsePKCS1PrivateKey(der); err == nil {
		return key, nil
	}

	if key, err = x509.ParsePKCS8PrivateKey(der); err == nil {
		switch key.(type) {
		case *rsa.PrivateKey, *ecdsa.PrivateKey:
			return
		default:
			return nil, fmt.Errorf("Found unknown private key type in PKCS#8 wrapping")
		}
	}

	if key, err = x509.ParseECPrivateKey(der); err == nil {
		return
	}

	return nil, fmt.Errorf("Invalid key type. The DER must contain an rsa.PrivateKey or ecdsa.PrivateKey")
}

func EncryptPEM(pemRaw []byte, passwd []byte) (*pem.Block, error) {
	block, _ := pem.Decode(pemRaw)
	if block == nil {
		return nil, fmt.Errorf("Failed decoding PEM. Block must be different from nil. [% x]", pemRaw)
	}

	der := block.Bytes

	privateKey, err := DERToPrivateKey(der)
	if err != nil {
		return nil, err
	}

	block, err = EncryptPrivateKey(privateKey, passwd)
	if err != nil {
		return nil, err
	}

	return block, nil
}

func EncryptPrivateKey(privateKey interface{}, passwd []byte) (*pem.Block, error) {
	switch k := privateKey.(type) {
	case *ecdsa.PrivateKey:
		raw, err := x509.MarshalECPrivateKey(k)
		if err != nil {
			return nil, err
		}

		block, err := x509.EncryptPEMBlock(rand.Reader, "EC PRIVATE KEY", raw, passwd, x509.PEMCipherAES256)
		if err != nil {
			return nil, err
		}

		return block, nil
	case *rsa.PrivateKey:
		raw := x509.MarshalPKCS1PrivateKey(k)

		block, err := x509.EncryptPEMBlock(rand.Reader, "RSA PRIVATE KEY", raw, passwd, x509.PEMCipherAES256)
		if err != nil {
			return nil, err
		}

		return block, nil
	default:
		return nil, fmt.Errorf("Invalid key type. It must be *ecdsa.PrivateKey or *rsa.PrivateKey")
	}
}

func TestAesCbC(t *testing.T) {

	oldStr := "fc4283c5-86ec-47ba-b181-4a1a6e4ca539"

	aaaa, err := CBCEncrypt(oldStr, "")

	fmt.Println(err)

	fmt.Println(aaaa)

	bbbb, _ := CBCDecrypt(aaaa, "")

	fmt.Println(bbbb)
}

var PubKey2 = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA1O3p0JN0/RrP7eY3f81i
zPf16FS0WMNGCJkd+y5c6yBzUvN0IEeoxiIWIBhoMKH0pzlzBg0rfttojSodOgNo
m/UCAzAYEgdIsNee5LSN/7e0T2/QvsIAHINuA8gI8fGoGiSA2TEzpUo6aVXwhZT3
4GGRdrSJ+m4iVk/Kt95tavBNk+NDVSeb5xAjxBchT5BjAMMlE0ffGZb0MMjjO5+e
9Tn8f99M2VMqpzXHXZzv1ABmqufzS20iWcSvnjhWcJ9hiKwO8Z30GgJyACmml+HM
xLYEFN9h2MWYgxLm9Z0rLMrWwMM+E2rCs8tsxAD5sO9RZMJPl1C0FIsMR53ngqbz
owIDAQAB
-----END PUBLIC KEY-----`

var PriKey = `-----BEGIN RSA PRIVATE KEY-----
MIIEpgIBAAKCAQEA1O3p0JN0/RrP7eY3f81izPf16FS0WMNGCJkd+y5c6yBzUvN0
IEeoxiIWIBhoMKH0pzlzBg0rfttojSodOgNom/UCAzAYEgdIsNee5LSN/7e0T2/Q
vsIAHINuA8gI8fGoGiSA2TEzpUo6aVXwhZT34GGRdrSJ+m4iVk/Kt95tavBNk+ND
VSeb5xAjxBchT5BjAMMlE0ffGZb0MMjjO5+e9Tn8f99M2VMqpzXHXZzv1ABmqufz
S20iWcSvnjhWcJ9hiKwO8Z30GgJyACmml+HMxLYEFN9h2MWYgxLm9Z0rLMrWwMM+
E2rCs8tsxAD5sO9RZMJPl1C0FIsMR53ngqbzowIDAQABAoIBAQCO1RE1ItUlO6kj
Un0ENAgEqojAUqGvsT33Yo7kAZO+/cOeb0UEqk0iq5bf7L9ncBynWDg6ZPc6X3/g
wdFdKxAvHck9zjM3VL+EMP+bNyrR0K8ZYk5Kx+Q/PEK+Mp8dfRdgggAUsZaNWB+a
rVVspiMo1wo28KBl5x8NevTnJkOLqXAyB7UyLWqnOL1fb988lZvZPR7ZUYroVIZa
pyXtZcafIJeKyQ3bvWI5+eFqOe61Z4Bx1+TpfZ3fKfSDW0vhxzNqaimOa8jSXtMJ
jMeOctL4nZ0TPo/jS3I+XlaH4ZQlFLuUWGscpxwfEeBN23I8HRLkZXJsw66yvRN3
s4bUKPXRAoGBAP/3oSZAECvfsYYzs76tnrAmR/0GxCqgguxDlWn5DowQzdWFOdHC
ZbTo/hUVoMSQnO1EKCFlnBS+wg/3TuIzUO0ewC1aeT7qHbOMDl0zKbNpS2Z9/j+U
zro+qz7XmkWolMCfmDrCrw9CtCxcMSII+ajbI8SAgFVMz9XnDt+xW9E9AoGBANT0
4F6kCUJTEyqf2+v84tjQ2wGIF6XtZPU9JR806zeMyahQ9F6z3hY8BYb0tIy5b3uJ
VlJ9TG1qg/t59TWxIq43mYSUJHe0aJi3ilooObQtHlhPu8nwmmX47sX0PyG2hMoD
kBVxTpTDmBaDz7O9uBnlMXJN5qEygctaixpEbmZfAoGBAMBA9kEMjRjnAyeRXcgy
D6aumhNqKZz6wltCx864yjxZwsBFOJBcOpgPCAg+HmqFU9jCAIJVF05dmNT1I8Ky
WG5BUoa+FaMzpOtenstRylh/Far9pyGKW1t4BpdEyRLY9CFZvbUk1OfZagqHlD/E
DgDN16eX/MwUzWYUDg/l3tjhAoGBAKGip/ZNjVWRFpggs9z/mfK1O7WC5Wgksp9N
ZLK2CN6l9p3RrFmBLk00C4HulGfHi+15RVLhFbRqx3iFje/N3iPbwaMWikNtZIKd
tN5Pb9To9gJTqpZRD+/cLOeFRrHBBjMK1z7fPKS/fN2B+JFVq7nD827t3+J0In4F
4FT0odMDAoGBAJk3ELB/FHY8xzZ4jF1wG/a1CK681Xm6SuU5KIELDSAUNoou6OPG
mS8gU20MMPAeV2z7khyDcSxlHsUyL73eLeaakbQov9NMW7cc99XX4wnP4W7FRpmr
QbHmKuHIRFHCFv+XX8c0aK2mDZMUlzJdy4FgD/YCEZ7kZMZKyvZW/ZuV
-----END RSA PRIVATE KEY-----`

func TestEncrypt(t *testing.T) {

	appSecret := "IgkibX71IEf382PT"
	encryptStr := "param_1=xxx&param_2=xxx&ak=xxx&ts=1111111111"

	snn, _ := Encrypt(encryptStr, []byte(appSecret), appSecret)
	fmt.Println(snn)

	// 验证签名
	pp, _ := Decrypt(snn, []byte(appSecret), appSecret)

	fmt.Println(pp)

	snn, _ = RSAPublicEncrypt(encryptStr, PubKey2)
	fmt.Println(snn)

	pp, _ = RSAPrivateDecrypt(snn, PriKey)
	fmt.Println(pp)

}

func TestBCrypt(t *testing.T) {

	appSecret := "中国人"
	//encryptStr := "param_1=xxx&param_2=xxx&ak=xxx&ts=1111111111"

	aaa, _ := BCryptPasswordEncoder(appSecret)

	fmt.Println(aaa)

	err := BCryptCompareHashAndPassword(aaa, appSecret)

	fmt.Println(err)

	bbb, _ := Argon2PasswordEncoder(appSecret)
	fmt.Println(bbb)

}
