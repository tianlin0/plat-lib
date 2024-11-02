package startupconfig

import "fmt"

// PaasAPI 服务接口抽象
type PaasAPI interface {
	// DomainName 接口域名
	DomainName() string
	// PolarisNamespace 北极星namespace
	PolarisNamespace() string
	// PolarisService 北极星service
	PolarisService() string
	// PolarisHost 北极星改写host
	PolarisHost() string
	// Url 接口Url
	Url(apiName string) string
	// AuthData 接口其他数据（鉴权数据等）
	AuthData(key string) (string, error)
	// UsePolaris 是否完整配置了北极星访问方式
	UsePolaris() bool
	// PolarisInstance 北极星结构实例
	PolarisInstance() *Polaris
}

// PaasApiConfig 服务接口
type PaasApiConfig struct {
	Domain  string               `json:"domain" yaml:"domain"`
	Polaris *Polaris             `json:"polaris" yaml:"polaris"`
	Auth    map[string]Decrypted `json:"auth" yaml:"auth"`
	Urls    map[string]string    `json:"urls" yaml:"urls"`
}

type Polaris struct {
	Host      string `json:"host" yaml:"host"`
	Namespace string `json:"namespace" yaml:"namespace"`
	Service   string `json:"service" yaml:"service"`
}

// DomainName 接口域名
func (c *PaasApiConfig) DomainName() string {
	if c == nil {
		return ""
	}
	return c.Domain
}

// PolarisNamespace 北极星namespace
func (c *PaasApiConfig) PolarisNamespace() string {
	if c == nil {
		return ""
	}
	if c.Polaris != nil {
		return c.Polaris.Namespace
	}
	return ""
}

// PolarisService 北极星service
func (c *PaasApiConfig) PolarisService() string {
	if c == nil {
		return ""
	}
	if c.Polaris != nil {
		return c.Polaris.Service
	}
	return ""
}

// PolarisHost 北极星改写host
func (c *PaasApiConfig) PolarisHost() string {
	if c == nil {
		return ""
	}
	if c.Polaris != nil {
		return c.Polaris.Host
	}
	return ""
}

// Url 接口Url
func (c *PaasApiConfig) Url(name string) string {
	if c == nil {
		return ""
	}
	if c.Urls != nil {
		return c.Urls[name]
	}
	return ""
}

// AuthData 接口其他数据（鉴权数据等）
func (c *PaasApiConfig) AuthData(key string) (string, error) {
	if c == nil {
		return "", fmt.Errorf("auth data %s empty", key)
	}
	if c.Auth != nil {
		if valueEncrypt, ok := c.Auth[key]; ok {
			return valueEncrypt.String(), nil
		}
	}
	return "", nil
}

// UsePolaris 接口访问是否配置了北极星
func (c *PaasApiConfig) UsePolaris() bool {
	if c == nil {
		return false
	}
	if polaris := c.Polaris; polaris != nil {
		if polaris.Host != "" && polaris.Service != "" && polaris.Namespace != "" {
			return true
		}
	}
	return false
}

// PolarisInstance 北极星结构实例
func (c *PaasApiConfig) PolarisInstance() *Polaris {
	if c == nil {
		return nil
	}
	return c.Polaris
}
