package startupconfig

import (
	"fmt"
)

// Database 数据库连接参数的抽象，包含使用sql.Open连接数据库时的参数
type Database interface {
	// DriverName 使用sql.Open连接数据库时的driverName参数
	DriverName() string
	// DatasourceName 使用sql.Open连接数据库时的datasourceName参数
	DatasourceName() string
	// ServerAddress 数据库服务器地址
	ServerAddress() string
	// Password 数据库用户密码
	Password() string
	// DatabaseName 数据库名称
	DatabaseName() interface{}
	// User 数据库用户
	User() string
	// Extend 扩展信息
	Extend(DatabaseExtendField) interface{}
}

// MysqlConfig mysql配置
type MysqlConfig struct {
	UserName          string    `json:"username" yaml:"username"`
	PasswordDecrypted Decrypted `json:"pwEncoded" yaml:"pwEncoded"`
	Address           string    `json:"address" yaml:"address"`
	Database          string    `json:"database" yaml:"database"`
}

// DriverName 使用sql.Open连接数据库时的driverName参数
func (c *MysqlConfig) DriverName() string {
	return "mysql"
}

// DatasourceName 使用sql.Open连接数据库时的datasourceName参数
func (c *MysqlConfig) DatasourceName() string {
	if c == nil {
		return ""
	}

	return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true&loc=Local",
		c.UserName, c.PasswordDecrypted, c.Address, c.Database)
}

// ServerAddress mysql服务器地址
func (c *MysqlConfig) ServerAddress() string {
	if c != nil {
		return c.Address
	}
	return ""
}

// Password mysql数据库用户密码
func (c *MysqlConfig) Password() string {
	if c == nil {
		return ""
	}
	return c.PasswordDecrypted.String()
}

// DatabaseName mysql数据库名称
func (c *MysqlConfig) DatabaseName() interface{} {
	if c == nil {
		return ""
	}
	return c.Database
}

// User mysql数据库用户
func (c *MysqlConfig) User() string {
	if c == nil {
		return ""
	}
	return c.UserName
}

// Extend 扩展字段
func (c *MysqlConfig) Extend(name DatabaseExtendField) interface{} {
	return nil
}

// RedisConfig redis配置
type RedisConfig struct {
	PasswordDecrypted Decrypted `json:"pwEncoded" yaml:"pwEncoded"`
	Address           string    `json:"address" yaml:"address"`
	Database          int64     `json:"database" yaml:"database"`
	Username          string    `json:"username" yaml:"username"`
	UseTLS            bool      `json:"useTLS" yaml:"useTLS"`
}

// DriverName 驱动名称
func (c *RedisConfig) DriverName() string {
	if c == nil {
		return ""
	}
	return "redis"
}

// DatasourceName 连接数据库时的datasourceName参数
func (c *RedisConfig) DatasourceName() string {
	return ""
}

// ServerAddress redis服务器地址
func (c *RedisConfig) ServerAddress() string {
	if c == nil {
		return ""
	}
	return c.Address
}

// Password redis数据库用户密码
func (c *RedisConfig) Password() string {
	if c == nil {
		return ""
	}
	return c.PasswordDecrypted.String()
}

// DatabaseName redis数据库名称
func (c *RedisConfig) DatabaseName() interface{} {
	if c == nil {
		return ""
	}
	return c.Database
}

// User redis数据库用户
func (c *RedisConfig) User() string {
	if c == nil {
		return ""
	}
	return c.Username
}

// Extend 扩展字段
func (c *RedisConfig) Extend(name DatabaseExtendField) interface{} {
	if c == nil {
		return ""
	}
	switch name {
	case redisUseTLS:
		return c.UseTLS
	}
	return nil
}

// DatabaseExtendField 数据库扩展字段名
type DatabaseExtendField string

var (
	redisUseTLS DatabaseExtendField = "TLS"
)
