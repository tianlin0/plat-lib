package startupconfig

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

// RunConfig 服务运行配置，有别于启动配置，运行配置是从配置文件或配置中心获取的配置，而非启动参数或环境变量获取的配置
type RunConfig interface {
	// MySQL 返回使用的MySQL连接参数
	MySQL(name string) Database
	// Redis 返回使用Redis连接参数
	Redis(name string) Database
	// TDMQ 返回使用的TDMQ连接参数
	TDMQ(name string) TDMQ
	// PaasAPI 服务接口参数
	PaasAPI(serviceName string) PaasAPI
	// Custom 自定义配置参数
	Custom() Custom
	// Trace trace配置参数
	Trace() Trace
	// Recorder 操作记录配置
	Recorder() Recorder
}

// StartupConfig 启动配置结构
type StartupConfig struct {
	MySQLMap            map[string]*MysqlConfig   `json:"mysql" yaml:"mysql"`
	RedisMap            map[string]*RedisConfig   `json:"redis" yaml:"redis"`
	TDMQMap             map[string]*TdmqConfig    `json:"tdmq" yaml:"tdmq"`
	ApiConfig           map[string]*PaasApiConfig `json:"api" yaml:"api"`
	CustomConfig        *CustomConfig             `json:"custom" yaml:"custom"`
	TracingConfig       *TracingConfig            `json:"tracing" yaml:"tracing"`
	UserOperationConfig *UserOperationRecorder    `json:"userOperation" yaml:"userOperation"`
}

// MySQL 返回使用的MySQL连接参数
func (c *StartupConfig) MySQL(name string) Database {
	if c != nil && c.MySQLMap != nil {
		if mysql, ok := c.MySQLMap[name]; ok {
			return mysql
		}
	}
	return nil
}

// Redis 返回使用Redis连接参数
func (c *StartupConfig) Redis(name string) Database {
	if c != nil && c.RedisMap != nil {
		if redis, ok := c.RedisMap[name]; ok {
			return redis
		}
	}
	return nil
}

// TDMQ 返回使用的TDMQ连接参数
func (c *StartupConfig) TDMQ(name string) TDMQ {
	if c != nil && c.TDMQMap != nil {
		if tdmq, ok := c.TDMQMap[name]; ok {
			return tdmq
		}
	}
	return nil
}

// PaasAPI 服务接口参数
func (c *StartupConfig) PaasAPI(serviceName string) PaasAPI {
	if c != nil && c.ApiConfig != nil {
		if api, ok := c.ApiConfig[serviceName]; ok {
			return api
		}
	}
	return nil
}

// Custom 返回自定义配置参数
func (c *StartupConfig) Custom() Custom {
	if c != nil {
		return c.CustomConfig
	}
	return nil
}

// Trace 返回Trace 配置参数
func (c *StartupConfig) Trace() Trace {
	if c != nil {
		return c.TracingConfig
	}
	return nil
}

// Recorder 操作记录配置
func (c *StartupConfig) Recorder() Recorder {
	if c != nil {
		return c.UserOperationConfig
	}
	return nil
}

// parseStartupConfig 从YAML中读取运行配置
func parseStartupConfig(rawConfig []byte) (*StartupConfig, error) {
	conf := &StartupConfig{}
	if err := yaml.Unmarshal(rawConfig, &conf); err != nil {
		return nil, err
	}
	return conf, nil
}

// newStartupConfig 初始化全局的运行参数，需要在启动参数初始化后再调用
func newStartupConfig(fileName string) (RunConfig, error) {
	configFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	runConfigInstance, err := parseStartupConfig(configFile)
	if err != nil {
		return nil, err
	}
	return runConfigInstance, nil
}
