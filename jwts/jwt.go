package jwts

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"time"
)

/*
user := jwt.Claims{
	Data: map[string]interface{}{
		"user_id": "mmmm",
	},
}
ppp, _ := jwt.SignToken(&user, "world", 6)
mmm, _ := jwt.ParseToken(ppp, "world")
*/

//Audience 受众
//IssuedAt 签发时间
//Issuer   发布人
//NotBefore 有效开始时间
//Subject  主题

// Claims ...
type Claims struct {
	Data map[string]interface{}
	jwt.StandardClaims
	SignMethod *jwt.SigningMethodHMAC
}

// SignToken 签发token
func SignToken(claims *Claims, secretKey string, seconds time.Duration) (sToken string, err error) {
	if claims == nil {
		return "", nil
	}

	//过期时间
	claims.ExpiresAt = time.Now().Add(seconds).Unix()
	//签署时间
	claims.IssuedAt = time.Now().Unix()

	if claims.SignMethod == nil {
		claims.SignMethod = jwt.SigningMethodHS256
	}

	token := jwt.NewWithClaims(claims.SignMethod, claims)
	return token.SignedString([]byte(secretKey))
}

// ParseToken 解析token
func ParseToken(sToken, secretKey string) (*Claims, error) {
	if sToken == "" || secretKey == "" {
		return nil, fmt.Errorf("ParseToken param is null")
	}
	var custom Claims
	token, err := jwt.ParseWithClaims(sToken, &custom, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected sign method %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, nil
}

// ParseTokenUnverified 直接返回内容
func ParseTokenUnverified(sToken string) (*Claims, error) {
	if sToken == "" {
		return nil, fmt.Errorf("ParseToken param is null")
	}

	p := jwt.Parser{}
	claims := &Claims{}
	if _, _, err := p.ParseUnverified(sToken, claims); err != nil {
		return nil, fmt.Errorf("unable to decode the access token: %v", err)
	}

	return claims, nil
}
