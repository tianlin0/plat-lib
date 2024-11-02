package startupconfig

// Custom 自定义配置抽象
type Custom interface {
	// GetSensitive 查询敏感配置（加密）对应key的value
	GetSensitive(key string) (string, error)
	// GetNormal 查询非敏感配置对应key的value
	GetNormal(key string) interface{}
}

type CustomConfig struct {
	Sensitive   map[string]Decrypted   `json:"sensitive" yaml:"sensitive"`
	Insensitive map[string]interface{} `json:"normal" yaml:"normal"`
}

// GetSensitive 查询敏感配置（加密）对应key的value
func (c *CustomConfig) GetSensitive(key string) (string, error) {
	if c != nil && c.Sensitive != nil {
		if valueEncrypted, ok := c.Sensitive[key]; ok {
			return valueEncrypted.String(), nil
		}
	}
	return "", nil
}

// GetNormal 查询非敏感配置对应key的value
func (c *CustomConfig) GetNormal(key string) interface{} {
	if c != nil && c.Insensitive != nil {
		if value, ok := c.Insensitive[key]; ok {
			return value
		}
	}
	return nil
}
