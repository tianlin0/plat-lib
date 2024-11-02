package startupconfig

import (
	"context"
	"fmt"
	"log"
	"os"
)

const (
	consulHost  = "CONSUL_HOST"
	consulToken = "CONSUL_TOKEN"
)

// ConsulHost 环境变量CONSUL_HOST
func ConsulHost() string {
	return os.Getenv(consulHost)
}

// ConsulToken 环境变量CONSUL_TOKEN
func ConsulToken() string {
	tokenEncrypted := os.Getenv(consulToken)
	token, err := DecDecrypt(tokenEncrypted)
	if err != nil {
		log.Fatal(context.TODO(), fmt.Sprintf(" decrypt CONSUL_TOKEN failed: %v", err))
		return ""
	}
	return token
}
