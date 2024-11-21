package jwts

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/tianlin0/plat-lib/conv"
	"testing"
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

func TestAesCbC(t *testing.T) {
	user := Claims{
		Data: map[string]interface{}{
			"user_id": "mmmm",
		},
		StandardClaims: jwt.StandardClaims{
			Subject: "aaaa",
		},
	}
	ppp, err1 := SignToken(&user, "world", 6*time.Second)

	//ppp = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJEYXRhIjp7InVzZXJfaWQiOiJtbW1tIn0sImV4cCI6MTcwNDQ0MTgwOCwiaWF0IjoxNzA0NDQxODAyLCJTaWduTWV0aG9kIjp7Ik5hbWUiOiJIUzI1NiIsIkhhc2giOjV9fQ.6yzrRetZ72kKlVJCAxtdPmLjbu4Yh_SYXi9Uazk61kA"
	mmm, err2 := ParseToken(ppp, "world")

	nnn, err3 := ParseTokenUnverified(ppp)

	fmt.Println(ppp, err1)
	fmt.Println(mmm, err2)
	fmt.Println(conv.String(nnn), err3)
}
