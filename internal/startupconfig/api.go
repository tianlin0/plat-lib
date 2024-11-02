package startupconfig

import (
	"log"
	"time"

	"github.com/json-iterator/go"
)

// ConfigAPI 配置访问实例
type ConfigAPI struct {
	runConfig   RunConfig
	configBytes []byte
	fileName    string
}

// New 创建一个配置实例
func New(fileName string) (*ConfigAPI, error) {
	conf, err := newStartupConfig(fileName)
	if err != nil {
		return nil, err
	}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	js, err := json.Marshal(conf)
	if err != nil {
		return nil, err
	}
	return &ConfigAPI{
		runConfig:   conf,
		fileName:    fileName,
		configBytes: js,
	}, nil
}

// NewByContent 通过配置内容创建实例
func NewByContent(content []byte) (*ConfigAPI, error) {
	conf, err := parseStartupConfig(content)
	if err != nil {
		return nil, err
	}
	return &ConfigAPI{
		runConfig: conf,
	}, nil
}

// Mysql Database interface
func (api *ConfigAPI) Mysql(name string) Database {
	if api.runConfig != nil {
		db := api.runConfig.MySQL(name)
		if db != nil {
			return db
		}
	}
	return nil
}

// MysqlDSN mysql Datasource Name of mysql name . eg: root:12345678@tcp(127.0.0.1:3306)/db_database?charset=utf8&parseTime=true&loc=Local
func (api *ConfigAPI) MysqlDSN(name string) string {
	if api.runConfig != nil {
		mySql := api.runConfig.MySQL(name)
		if mySql != nil {
			return mySql.DatasourceName()
		}
	}
	return ""
}

// MysqlDriverName mysql Driver Name of mysql name
func (api *ConfigAPI) MysqlDriverName(name string) string {
	if api.runConfig != nil {
		mySql := api.runConfig.MySQL(name)
		if mySql != nil {
			return mySql.DriverName()
		}
	}
	return ""
}

// MysqlPassword mysql password of mysql name
func (api *ConfigAPI) MysqlPassword(name string) string {
	if api.runConfig != nil {
		mySql := api.runConfig.MySQL(name)
		if mySql != nil {
			return mySql.Password()
		}
	}
	return ""
}

// MysqlAddress mysql server address of mysql name
func (api *ConfigAPI) MysqlAddress(name string) string {
	if api.runConfig != nil {
		mySql := api.runConfig.MySQL(name)
		if mySql != nil {
			return mySql.ServerAddress()
		}
	}
	return ""
}

// MysqlUser mysql user of mysql name
func (api *ConfigAPI) MysqlUser(name string) string {
	if api.runConfig != nil {
		mySql := api.runConfig.MySQL(name)
		if mySql != nil {
			return mySql.User()
		}
	}
	return ""
}

// MysqlDatabase mysql database of mysql name
func (api *ConfigAPI) MysqlDatabase(name string) string {
	if api.runConfig != nil {
		mySql := api.runConfig.MySQL(name)
		if mySql != nil {
			name, ok := mySql.DatabaseName().(string)
			if !ok {
				return ""
			}
			return name
		}
	}
	return ""
}

// Redis Database interface
func (api *ConfigAPI) Redis(name string) Database {
	if api.runConfig != nil {
		db := api.runConfig.Redis(name)
		if db != nil {
			return db
		}
	}
	return nil
}

// RedisAddress redis Address of redis name
func (api *ConfigAPI) RedisAddress(name string) string {
	if api.runConfig != nil {
		redis := api.runConfig.Redis(name)
		if redis != nil {
			return redis.ServerAddress()
		}
	}
	return ""
}

// RedisPWD redis password of redis name
func (api *ConfigAPI) RedisPWD(name string) string {
	if api.runConfig != nil {
		redis := api.runConfig.Redis(name)
		if redis != nil {
			return redis.Password()
		}
	}
	return ""
}

// RedisDatabase redis password of redis name
func (api *ConfigAPI) RedisDatabase(name string) int64 {
	if api.runConfig != nil {
		redis := api.runConfig.Redis(name)
		if redis != nil {
			db, ok := redis.DatabaseName().(int64)
			if !ok {
				return 0
			}
			return db
		}
	}
	return 0
}

// RedisUser redis user of redis name
func (api *ConfigAPI) RedisUser(name string) string {
	if api.runConfig != nil {
		redis := api.runConfig.Redis(name)
		if redis != nil {
			return redis.User()
		}
	}
	return ""
}

// RedisUseTLS redis use TLS or not of redis name
func (api *ConfigAPI) RedisUseTLS(name string) bool {
	if api.runConfig != nil {
		redis := api.runConfig.Redis(name)
		if redis != nil {
			use, _ := redis.Extend(redisUseTLS).(bool)
			return use
		}
	}
	return false
}

// TDMQ  interface
func (api *ConfigAPI) TDMQ(name string) TDMQ {
	if api.runConfig != nil {
		mq := api.runConfig.TDMQ(name)
		if mq != nil {
			return mq
		}
	}
	return nil
}

// TDMQUrl TDMQ URL of tdmq name
func (api *ConfigAPI) TDMQUrl(name string) string {
	if api.runConfig != nil {
		mq := api.runConfig.TDMQ(name)
		if mq != nil {
			return mq.URL()
		}
	}
	return ""
}

// TDMQSubscription TDMQ Subscription Name of tdmq name
func (api *ConfigAPI) TDMQSubscription(name string) string {
	if api.runConfig != nil {
		mq := api.runConfig.TDMQ(name)
		if mq != nil {
			return mq.SubscriptionName()
		}
	}
	return ""
}

// TDMQToken TDMQ token of tdmq name
func (api *ConfigAPI) TDMQToken(name string) string {
	if api.runConfig != nil {
		mq := api.runConfig.TDMQ(name)
		if mq != nil {
			return mq.AuthenticationToken()
		}
	}
	return ""
}

// TDMQInitialPosition TDMQ Subscription Initial Position  of tdmq name
func (api *ConfigAPI) TDMQInitialPosition(name string) int {
	if api.runConfig != nil {
		mq := api.runConfig.TDMQ(name)
		if mq != nil {
			return mq.SubscriptionInitialPosition()
		}
	}
	return SubscriptionPositionLatest
}

// TDMQTopic TDMQ topic of name of tdmq name
func (api *ConfigAPI) TDMQTopic(name, topicName string) string {
	if api.runConfig != nil {
		mq := api.runConfig.TDMQ(name)
		if mq != nil {
			return mq.Topic(topicName)
		}
	}
	return ""
}

// PaasAPI PaasAPI interface
func (api *ConfigAPI) PaasAPI(service string) PaasAPI {
	if api.runConfig != nil {
		passApi := api.runConfig.PaasAPI(service)
		if passApi != nil {
			return passApi
		}
	}
	return nil
}

// UserPolaris can use polaris
func (api *ConfigAPI) UserPolaris(service string) bool {
	if api.runConfig != nil {
		srv := api.runConfig.PaasAPI(service)
		if srv != nil {
			return srv.UsePolaris()
		}
	}
	return false
}

// PaasAPIDomain paas api domain name
func (api *ConfigAPI) PaasAPIDomain(service string) string {
	if api.runConfig != nil {
		srv := api.runConfig.PaasAPI(service)
		if srv != nil {
			return srv.DomainName()
		}
	}
	return ""
}

// PaasAPIPolaris paas api polaris instance. It is recommended not to modify
func (api *ConfigAPI) PaasAPIPolaris(service string) *Polaris {
	if api.runConfig != nil {
		srv := api.runConfig.PaasAPI(service)
		if srv != nil {
			return srv.PolarisInstance()
		}
	}
	return nil
}

// PaasAPIPolarisHost paas api polaris host
func (api *ConfigAPI) PaasAPIPolarisHost(service string) string {
	if api.runConfig != nil {
		srv := api.runConfig.PaasAPI(service)
		if srv != nil {
			return srv.PolarisHost()
		}
	}
	return ""
}

// PaasAPIPolarisNamespace paas api polaris namespace
func (api *ConfigAPI) PaasAPIPolarisNamespace(service string) string {
	if api.runConfig != nil {
		srv := api.runConfig.PaasAPI(service)
		if srv != nil {
			return srv.PolarisNamespace()
		}
	}
	return ""
}

// PaasAPIPolarisService paas api polaris Service
func (api *ConfigAPI) PaasAPIPolarisService(service string) string {
	if api.runConfig != nil {
		srv := api.runConfig.PaasAPI(service)
		if srv != nil {
			return srv.PolarisService()
		}
	}
	return ""
}

// PaasAPIUrl paas api url
func (api *ConfigAPI) PaasAPIUrl(service, apiName string) string {
	if api.runConfig != nil {
		srv := api.runConfig.PaasAPI(service)
		if srv != nil {
			return srv.Url(apiName)
		}
	}
	return ""
}

// PaasAPIAuthValue paas api auth value(encrypted) of key
func (api *ConfigAPI) PaasAPIAuthValue(service, key string) (string, error) {
	if api.runConfig != nil {
		srv := api.runConfig.PaasAPI(service)
		if srv != nil {
			return srv.AuthData(key)
		}
	}
	return "", nil
}

// Custom custom config interface
func (api *ConfigAPI) Custom() Custom {
	if api.runConfig != nil {
		cs := api.runConfig.Custom()
		if cs != nil {
			return cs
		}
	}
	return nil
}

// CustomSensitive get value of sensitive custom config key (value encrypted)
func (api *ConfigAPI) CustomSensitive(key string) (string, error) {
	if api.runConfig != nil {
		conf := api.runConfig.Custom()
		if conf != nil {
			value, err := conf.GetSensitive(key)
			if err != nil {
				return "", err
			}
			return value, nil
		}
	}
	return "", nil
}

// CustomNormal get value of insensitive custom config key
func (api *ConfigAPI) CustomNormal(key string) interface{} {
	if api.runConfig != nil {
		conf := api.runConfig.Custom()
		if conf != nil {
			normal := conf.GetNormal(key)
			if normal != nil {
				return normal
			}
		}
	}
	return nil
}

// ConsulHost get consul host from env
func (api *ConfigAPI) ConsulHost() string {
	return ConsulHost()
}

// ConsulToken get consul token from env
func (api *ConfigAPI) ConsulToken() string {
	return ConsulToken()
}

// Trace get Trace interface
func (api *ConfigAPI) Trace() Trace {
	if api.runConfig != nil {
		trace := api.runConfig.Trace()
		if trace != nil {
			return trace
		}
	}
	return nil
}

// TraceConfig get TracingConfig Instance
func (api *ConfigAPI) TraceConfig() *TracingConfig {
	if conf, ok := api.Trace().(*TracingConfig); ok {
		return conf
	}
	return nil
}

// Recorder get Recorder instance
func (api *ConfigAPI) Recorder() Recorder {
	if api.runConfig != nil {
		rd := api.runConfig.Recorder()
		if rd != nil {
			return rd
		}
	}
	return nil
}

// RecorderClosed if close the recorder
// @receiver api
// @return bool
func (api *ConfigAPI) RecorderClosed() bool {
	rd := api.Recorder()
	if rd == nil {
		return false
	}
	return rd.Close()
}

// RecorderTDMQ get Recorder TDMQ configuration instance
func (api *ConfigAPI) RecorderTDMQ() TDMQ {
	rd := api.Recorder()
	if rd == nil {
		return nil
	}
	key := rd.MQConfigKey()
	if key == "" {
		return nil
	}
	mq := api.TDMQ(key)
	if mq != nil {
		return mq
	}
	return nil
}

// RecorderTDMQURL get Recorder TDMQ URL
func (api *ConfigAPI) RecorderTDMQURL() string {
	mq := api.RecorderTDMQ()
	if mq != nil {
		return mq.URL()
	}
	return ""
}

// RecorderTDMQToken get Recorder TDMQ token
func (api *ConfigAPI) RecorderTDMQToken() string {
	mq := api.RecorderTDMQ()
	if mq != nil {
		return mq.AuthenticationToken()
	}
	return ""
}

// RecorderTDMQTopic get Recorder TDMQ topic
func (api *ConfigAPI) RecorderTDMQTopic() string {
	mq := api.RecorderTDMQ()
	if mq != nil {
		return mq.Topic(userOperationTopicKey)
	}
	return ""
}

// RecorderGroups get Recorder groups information
func (api *ConfigAPI) RecorderGroups() map[string]*OperationGroup {
	rd := api.Recorder()
	if rd == nil {
		return nil
	}
	groups := rd.OperationGroups()
	if groups != nil {
		return groups
	}
	return nil
}

// RecorderGroup get Recorder group of key
// @receiver api
// @param key is group key
// @return *OperationGroup
func (api *ConfigAPI) RecorderGroup(key string) *OperationGroup {
	rd := api.Recorder()
	if rd == nil {
		return nil
	}
	group := rd.OperationGroup(key)
	if group != nil {
		return group
	}
	return nil
}

// RunConfig RunConfig interface
func (api *ConfigAPI) RunConfig() RunConfig {
	if api.runConfig != nil {
		return api.runConfig
	}
	return nil
}

// StartupConfig get StartupConfig Instance
func (api *ConfigAPI) StartupConfig() *StartupConfig {
	if conf, ok := api.runConfig.(*StartupConfig); ok {
		return conf
	}
	return nil
}

// StartupMysqlAll get all mysql Instances
func (api *ConfigAPI) StartupMysqlAll() map[string]*MysqlConfig {
	if conf, ok := api.runConfig.(*StartupConfig); ok {
		return conf.MySQLMap
	}
	return nil
}

// StartupRedisAll get all redis Instances
func (api *ConfigAPI) StartupRedisAll() map[string]*RedisConfig {
	if conf, ok := api.runConfig.(*StartupConfig); ok {
		return conf.RedisMap
	}
	return nil
}

// StartupTDMQAll get all tdmq Instances
func (api *ConfigAPI) StartupTDMQAll() map[string]*TdmqConfig {
	if conf, ok := api.runConfig.(*StartupConfig); ok {
		return conf.TDMQMap
	}
	return nil
}

// StartupCustomSensitiveAll get all custom sensitive configs(kv)
func (api *ConfigAPI) StartupCustomSensitiveAll() map[string]Decrypted {
	if conf, ok := api.runConfig.(*StartupConfig); ok {
		if conf.CustomConfig != nil {
			return conf.CustomConfig.Sensitive
		}
	}
	return nil
}

// StartupCustomNormalAll get all custom normal configs(kv)
func (api *ConfigAPI) StartupCustomNormalAll() map[string]interface{} {
	if conf, ok := api.runConfig.(*StartupConfig); ok {
		if conf.CustomConfig != nil {
			return conf.CustomConfig.Insensitive
		}
	}
	return nil
}

// StartupPaasApiAll get all api configs
func (api *ConfigAPI) StartupPaasApiAll() map[string]*PaasApiConfig {
	if conf, ok := api.runConfig.(*StartupConfig); ok {
		if conf.ApiConfig != nil {
			return conf.ApiConfig
		}
	}
	return nil
}

// Transform 将配置内容转换为指定结构体, 加密字段会自动解密
func (api *ConfigAPI) Transform(to interface{}) error {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	if err := json.Unmarshal(api.configBytes, to); err != nil {
		return err
	}
	return nil
}

// PollUpdate poll for updates, default update per 60s
func (api *ConfigAPI) PollUpdate(callback func(api *ConfigAPI) error, duration ...time.Duration) {
	dur := 60 * time.Second
	if len(duration) > 0 {
		dur = duration[0]
	}
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Println("poll for update panic: ", err)
			}
		}()
		for {
			t := time.NewTimer(dur)
			<-t.C
			conf, err := newStartupConfig(api.fileName)
			if err != nil {
				log.Println("Failed to poll for update config: ", err)
				continue
			}
			api.runConfig = conf
			if err := callback(api); err != nil {
				log.Println("Callback failed: ", err)
			}
		}
	}()
}

// Update update config from file
func (api *ConfigAPI) Update() error {
	conf, err := newStartupConfig(api.fileName)
	if err != nil {
		return err
	}
	api.runConfig = conf
	return nil
}
