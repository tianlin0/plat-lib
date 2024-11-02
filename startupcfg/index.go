package startupcfg

import (
	"fmt"
	"git.woa.com/odp-go/gdp-lib/conn"
	"git.woa.com/odp-go/gdp-lib/conv"
	"git.woa.com/odp-go/gdp-lib/logs"
	startupCfg "git.woa.com/odp-go/gdp-lib/startupcfg/startupconfig"
	"git.woa.com/odp-go/gdp-lib/templates"
	"net"
	"net/url"
	"strings"
)

const (
	envDefault      = "default" //如果没有设置环境变量，则使用默认值
	envSplitString  = "/"
	paasSplitString = "_$_$_$_"
)

func getInstanceFromYaml(isFilePathName bool, configFile string) (*startupCfg.ConfigAPI, error) {
	var conf *startupCfg.ConfigAPI
	var err error

	if isFilePathName {
		conf, err = startupCfg.New(configFile)
	} else {
		conf, err = startupCfg.NewByContent([]byte(configFile))
	}

	if err != nil {
		logs.DefaultLogger().Error(err)
		return nil, err
	}
	return conf, nil
}

func getEnvAndKeyName(keyString string, envDefault string, envSplitString string) (string, string) {
	keyList := strings.Split(keyString, envSplitString)
	envTemp := envDefault
	keyNameTemp := ""

	if len(keyList) >= 2 { //从后往前获取keyName和envName
		lenInt := len(keyList)
		envTemp = strings.TrimSpace(keyList[lenInt-2])
		keyNameTemp = strings.TrimSpace(keyList[lenInt-1])
	} else if len(keyList) == 1 {
		keyNameTemp = strings.TrimSpace(keyList[0])
	}

	if envTemp == "" {
		envTemp = envDefault
	}

	return envTemp, keyNameTemp
}

func getOneConnFromAddress(address string) (*conn.Connect, error) {
	oneConnect := new(conn.Connect)

	//如果包含http
	hostPort := ""

	urlTemp, err := url.Parse(address)
	if err == nil {
		if urlTemp.Scheme != "" {
			oneConnect.Protocol = urlTemp.Scheme
		}
		if urlTemp.Host != "" {
			hostPort = urlTemp.Host
		}
	}

	if hostPort == "" {
		hostPort = address
	}

	host, port, err := net.SplitHostPort(hostPort)
	if err != nil {
		logs.DefaultLogger().Error("getOneConnFromAddress Error:", address, err)
		return nil, err
	}
	oneConnect.Host = host
	oneConnect.Port = port
	return oneConnect, nil
}

func (s *Startup) getEnvSplitKeyName(env, keyName string) string {
	if env == "" {
		env = s.envDefault
	}
	if keyName == "" {
		return ""
	}
	return fmt.Sprintf("%s%s%s", env, s.envSplitString, keyName)
}

func getConnFromMysql(mysqlCfg *startupCfg.MysqlConfig) *conn.Connect {
	oneConnect, err := getOneConnFromAddress(mysqlCfg.ServerAddress())
	if err != nil {
		return nil
	}

	oneConnect.Driver = conn.DriverMysql
	oneConnect.Username = mysqlCfg.User()
	oneConnect.Password = mysqlCfg.Password()
	oneConnect.Database = conv.String(mysqlCfg.DatabaseName())

	return oneConnect
}
func getConnFromRedis(redisCfg *startupCfg.RedisConfig, useTLS bool) *conn.Connect {
	oneConnect, err := getOneConnFromAddress(redisCfg.ServerAddress())
	if err != nil {
		return nil
	}

	oneConnect.Driver = conn.DriverRedis
	oneConnect.Username = redisCfg.User()
	oneConnect.Password = redisCfg.Password()
	oneConnect.Database = conv.String(redisCfg.DatabaseName())
	if useTLS {
		oneConnect.Extend = map[string]interface{}{
			"useTLS": useTLS,
		}
	}

	return oneConnect
}
func getConnFromTdmq(tdmqCfg *startupCfg.TdmqConfig) *conn.Connect {
	oneConnect, err := getOneConnFromAddress(tdmqCfg.URL())
	if err != nil {
		return nil
	}

	oneConnect.Driver = conn.DriverTdmq
	oneConnect.Password = tdmqCfg.AuthenticationToken()
	oneConnect.Extend = map[string]interface{}{}
	{ //处理配置的topics列表
		oneConnect.Extend["topics"] = tdmqCfg.Topics
	}

	return oneConnect
}

// env/keyName
func (s *Startup) addConnToMap(connectMap map[string]*conn.Connect, envKeyName string, con *conn.Connect) map[string]*conn.Connect {
	env, keyName := getEnvAndKeyName(envKeyName, s.envDefault, s.envSplitString)
	if env != "" && keyName != "" {
		mapKeyName := s.getEnvSplitKeyName(env, keyName)
		connectMap[mapKeyName] = con
	}
	return connectMap
}

// 优化构建url列表
// env/paasName/keyName
func (s *Startup) getAllUrlMaps(apiMap map[string]*startupCfg.PaasApiConfig) map[string]*urlStruct {
	//取出同一个服务的所有url的key。paasName/keyName = url
	//后面的不覆盖前面的
	allUrlMap := make(map[string]map[string]string)
	for key, val := range apiMap {
		_, paasName := getEnvAndKeyName(key, s.envDefault, s.envSplitString)
		if _, ok := allUrlMap[paasName]; !ok {
			allUrlMap[paasName] = make(map[string]string)
		}
		for keyName, url := range val.Urls {
			if _, ok := allUrlMap[paasName][keyName]; !ok {
				allUrlMap[paasName][keyName] = url
			}
		}
	}

	allApiCfgMap := make(map[string]*urlStruct)

	//先填充全局默认的，然后填充默认的
	for key, val := range apiMap {
		//Url进行扩容，因为如果是同一个服务，有可能不会每个地址都重复写一遍,或者也有遗漏的情况
		envTemp, paasName := getEnvAndKeyName(key, s.envDefault, s.envSplitString)
		{
			if urlMap, ok := allUrlMap[paasName]; ok {
				if val.Urls == nil {
					val.Urls = make(map[string]string)
				}
				for gdpKey, url := range urlMap {
					//如果不存在，则用默认的，如果特殊设置了，则用此特殊设置的
					if _, ok1 := val.Urls[gdpKey]; !ok1 {
						val.Urls[gdpKey] = url
					}
				}
			}
		}

		if len(val.Urls) == 0 {
			continue
		}

		for keyName, urlVal := range val.Urls {
			mapKey := fmt.Sprintf("%s%s%s%s%s", envTemp, s.envSplitString, paasName, s.paasSplitString, keyName)
			urlCfg := new(urlStruct)
			urlCfg.PaasName = paasName
			urlCfg.ApiName = keyName
			urlCfg.PaasAPIUrl = val.DomainName() + urlVal
			urlCfg.paasAPICfg = val
			allApiCfgMap[mapKey] = urlCfg
		}
	}
	return allApiCfgMap
}

func (s *Startup) getAllCustom(customSenMap map[string]startupCfg.Decrypted, customMap map[string]interface{}) (map[string]interface{}, error) {
	secretMap := make(map[string]string)
	for key, val := range customSenMap {
		oldStr := val.String()
		secretMap[key] = oldStr
	}

	customMapNew := customMap
	if len(secretMap) > 0 {
		customStr := conv.String(customMap)
		if customStr != "" {
			postData, err := templates.Template(customStr, secretMap)
			if err != nil {
				return nil, err
			}
			customMapNew = make(map[string]interface{})
			err = conv.Unmarshal(postData, &customMapNew)
			if err != nil {
				return nil, err
			}
		}
	}

	allCustomMap := make(map[string]interface{})

	for key, val := range customMapNew {
		envTemp, keyNameTemp := getEnvAndKeyName(key, s.envDefault, s.envSplitString)
		if envTemp != "" && keyNameTemp != "" {
			mapKeyName := s.getEnvSplitKeyName(envTemp, keyNameTemp)
			allCustomMap[mapKeyName] = val
		}
	}

	return allCustomMap, nil
}

func getMysqlFromYaml(conf *startupCfg.ConfigAPI) map[string]*startupCfg.MysqlConfig {
	mysqlMap := conf.StartupMysqlAll()
	mysqlAllMap := make(map[string]*startupCfg.MysqlConfig)
	for key, val := range mysqlMap {
		if key != "" && val != nil {
			mysqlAllMap[key] = val
		}
	}
	return mysqlAllMap
}
func getRedisFromYaml(conf *startupCfg.ConfigAPI) map[string]*startupCfg.RedisConfig {
	redisMap := conf.StartupRedisAll()
	redisAllMap := make(map[string]*startupCfg.RedisConfig)
	for key, val := range redisMap {
		if key != "" && val != nil {
			redisAllMap[key] = val
		}
	}
	return redisAllMap
}
func getTdmqFromYaml(conf *startupCfg.ConfigAPI) map[string]*startupCfg.TdmqConfig {
	tdmqMap := conf.StartupTDMQAll()
	tdmqAllMap := make(map[string]*startupCfg.TdmqConfig)
	for key, val := range tdmqMap {
		if key != "" && val != nil {
			tdmqAllMap[key] = val
		}
	}
	return tdmqAllMap
}

func (s *Startup) getAllMysqlFromYaml(conf *startupCfg.ConfigAPI) map[string]*startupCfg.MysqlConfig {
	allMysqlMap := make(map[string]*startupCfg.MysqlConfig)

	mysqlMap := getMysqlFromYaml(conf)
	//先填充全局默认的，然后填充默认的
	for key, val := range mysqlMap {
		envTemp, keyNameTemp := getEnvAndKeyName(key, s.envDefault, s.envSplitString)
		if envTemp != "" && keyNameTemp != "" {
			mapKeyName := s.getEnvSplitKeyName(envTemp, keyNameTemp)
			allMysqlMap[mapKeyName] = val
		}
	}
	return allMysqlMap
}
func (s *Startup) getAllRedisFromYaml(conf *startupCfg.ConfigAPI) map[string]*startupCfg.RedisConfig {
	allRedisMap := make(map[string]*startupCfg.RedisConfig)

	redisMap := getRedisFromYaml(conf)
	//先填充全局默认的，然后填充默认的
	for key, val := range redisMap {
		envTemp, keyNameTemp := getEnvAndKeyName(key, s.envDefault, s.envSplitString)
		if envTemp != "" && keyNameTemp != "" {
			mapKeyName := s.getEnvSplitKeyName(envTemp, keyNameTemp)
			allRedisMap[mapKeyName] = val
		}
	}
	return allRedisMap
}
func (s *Startup) getAllTdmqFromYaml(conf *startupCfg.ConfigAPI) map[string]*startupCfg.TdmqConfig {
	allTdmqMap := make(map[string]*startupCfg.TdmqConfig)

	tdmqMap := getTdmqFromYaml(conf)
	//先填充全局默认的，然后填充默认的
	for key, val := range tdmqMap {
		envTemp, keyNameTemp := getEnvAndKeyName(key, s.envDefault, s.envSplitString)
		if envTemp != "" && keyNameTemp != "" {
			mapKeyName := s.getEnvSplitKeyName(envTemp, keyNameTemp)
			allTdmqMap[mapKeyName] = val
		}
	}
	return allTdmqMap
}

// getAllConnectFromYaml 返回的是某环境下的连接 key的格式必须为:环境/变量名，比如：loc/MysqlConnect  返回的对象是，环境，变量名，连接
func (s *Startup) getAllConnectFromYaml(conf *startupCfg.ConfigAPI) map[string]*conn.Connect {
	connectMap := make(map[string]*conn.Connect)

	{ //获取mysql列表，环境，然后该环境下所有变量名
		mysqlMap := getMysqlFromYaml(conf)
		//先填充全局默认的，然后填充默认的
		for key, val := range mysqlMap {
			oneConnect := getConnFromMysql(val)
			if oneConnect == nil {
				continue
			}
			connectMap = s.addConnToMap(connectMap, key, oneConnect)
		}
	}

	{ //获取redis列表
		redisMap := getRedisFromYaml(conf)
		for key, val := range redisMap {
			oneConnect := getConnFromRedis(val, conf.RedisUseTLS(key))
			if oneConnect == nil {
				continue
			}
			connectMap = s.addConnToMap(connectMap, key, oneConnect)
		}
	}

	{ //获取tdmq列表
		tdmqMap := getTdmqFromYaml(conf)
		for key, val := range tdmqMap {
			oneConnect := getConnFromTdmq(val)
			if oneConnect == nil {
				continue
			}
			connectMap = s.addConnToMap(connectMap, key, oneConnect)
		}
	}

	return connectMap
}

// getAllApiUrlFromYaml 获取某环境下的所有地址  env.paasName.keyName
func (s *Startup) getAllApiUrlFromYaml(conf *startupCfg.ConfigAPI) map[string]*urlStruct {
	apiMap := conf.StartupPaasApiAll()
	return s.getAllUrlMaps(apiMap)
}

// getAllCustomFromYaml 加密的数据需要用""引起来，避免出错，且key不能包含-/等特殊字符
func (s *Startup) getAllCustomFromYaml(conf *startupCfg.ConfigAPI) (map[string]interface{}, error) {
	customSenMap := conf.StartupCustomSensitiveAll()
	customMap := conf.StartupCustomNormalAll()
	return s.getAllCustom(customSenMap, customMap)
}

type urlStruct struct {
	PaasName   string
	ApiName    string
	paasAPICfg *startupCfg.PaasApiConfig
	PaasAPIUrl string
}

// GetPaasApiCfg 直接获取url全地址
func (u *urlStruct) GetPaasApiCfg() *startupCfg.PaasApiConfig {
	return u.paasAPICfg
}

// Startup 初始化一个自定义配置
type Startup struct {
	Mysql   map[string]*startupCfg.MysqlConfig
	Redis   map[string]*startupCfg.RedisConfig
	Tdmq    map[string]*startupCfg.TdmqConfig
	Connect map[string]*conn.Connect
	Api     map[string]*urlStruct
	Custom  map[string]interface{}

	envDefault      string //没有配置env的默认值
	envSplitString  string //envSplitString 与 paasSplitString 不能相同
	paasSplitString string
}

type cfgOpt func(*Startup)

// NewOption 新增
func NewOption() cfgOpt {
	return func(*Startup) {}
}

// SetEmptyEnvKeyName 设置当env为空时的默认key
func (c cfgOpt) SetEmptyEnvKeyName(driver string) cfgOpt {
	return func(do *Startup) {
		c(do)
		do.envDefault = driver
	}
}

//SetEnvSplit env的分隔符
func (c cfgOpt) SetEnvSplit(split string) cfgOpt {
	return func(do *Startup) {
		c(do)
		if split != do.paasSplitString {
			do.envSplitString = split
		}
	}
}

//SetPaasSplit url中paas与urlKey之间的分隔符
func (c cfgOpt) SetPaasSplit(split string) cfgOpt {
	return func(do *Startup) {
		c(do)
		if split != do.envSplitString {
			do.paasSplitString = split
		}
	}
}

func commGetStartup(conf *startupCfg.ConfigAPI, options ...cfgOpt) (*Startup, error) {
	startTemp := new(Startup)

	//设置变量
	for _, oneFun := range options {
		oneFun(startTemp)
	}

	//没有设置，就采用默认值
	if startTemp.envDefault == "" {
		startTemp.envDefault = envDefault
	}
	if startTemp.envSplitString == "" {
		startTemp.envSplitString = envSplitString
	}
	if startTemp.paasSplitString == "" {
		startTemp.paasSplitString = paasSplitString
	}

	startTemp.Mysql = startTemp.getAllMysqlFromYaml(conf)
	startTemp.Redis = startTemp.getAllRedisFromYaml(conf)
	startTemp.Tdmq = startTemp.getAllTdmqFromYaml(conf)

	startTemp.Connect = startTemp.getAllConnectFromYaml(conf)
	startTemp.Api = startTemp.getAllApiUrlFromYaml(conf)

	customTemp, err := startTemp.getAllCustomFromYaml(conf)
	if err != nil {
		return nil, err
	}
	startTemp.Custom = customTemp

	return startTemp, nil
}

// NewStartupForYamlFile 初始化一个yaml文件的Startup配置
func NewStartupForYamlFile(configFile string, options ...cfgOpt) (*Startup, error) {
	conf, err := getInstanceFromYaml(true, configFile)
	if err != nil {
		return nil, err
	}
	return commGetStartup(conf, options...)
}

// NewStartupForYamlContent 初始化一个yaml内容的Startup配置
func NewStartupForYamlContent(configContent string, options ...cfgOpt) (*Startup, error) {
	conf, err := getInstanceFromYaml(false, configContent)
	if err != nil {
		return nil, err
	}
	return commGetStartup(conf, options...)
}

func (s *Startup) getKeyListByEnv(env string, keyList []string) map[string]string {
	defaultPrefix := s.envDefault + s.envSplitString
	envPrefix := env + s.envSplitString

	newKeyMap := map[string]string{}
	for _, one := range keyList {
		if strings.HasPrefix(one, defaultPrefix) {
			keyName := strings.TrimPrefix(one, defaultPrefix)
			newKeyMap[keyName] = one
		}
	}
	for _, one := range keyList {
		if strings.HasPrefix(one, envPrefix) {
			keyName := strings.TrimPrefix(one, envPrefix)
			newKeyMap[keyName] = one
		}
	}
	return newKeyMap
}

// getAllMysql 取得某一个环境的所有变量
func (s *Startup) getAllMysql(env string) map[string]*startupCfg.MysqlConfig {
	allConnect := make(map[string]*startupCfg.MysqlConfig)

	keyList := make([]string, 0)
	for key, _ := range s.Mysql {
		keyList = append(keyList, key)
	}
	keyMap := s.getKeyListByEnv(env, keyList)
	for keyName, oldKeyName := range keyMap {
		allConnect[keyName] = s.Mysql[oldKeyName]
	}
	return allConnect
}

// GetAllRedis 取得某一个环境的所有变量
func (s *Startup) getAllRedis(env string) map[string]*startupCfg.RedisConfig {
	allConnect := make(map[string]*startupCfg.RedisConfig)

	keyList := make([]string, 0)
	for key, _ := range s.Redis {
		keyList = append(keyList, key)
	}
	keyMap := s.getKeyListByEnv(env, keyList)
	for keyName, oldKeyName := range keyMap {
		allConnect[keyName] = s.Redis[oldKeyName]
	}
	return allConnect
}

// GetAllTdmq 取得某一个环境的所有变量
func (s *Startup) getAllTdmq(env string) map[string]*startupCfg.TdmqConfig {
	allConnect := make(map[string]*startupCfg.TdmqConfig)

	keyList := make([]string, 0)
	for key, _ := range s.Tdmq {
		keyList = append(keyList, key)
	}
	keyMap := s.getKeyListByEnv(env, keyList)
	for keyName, oldKeyName := range keyMap {
		allConnect[keyName] = s.Tdmq[oldKeyName]
	}
	return allConnect
}

// GetAllConnect 取得某一个环境的所有变量
func (s *Startup) getAllConnect(env string) map[string]*conn.Connect {
	allConnect := make(map[string]*conn.Connect)

	keyList := make([]string, 0)
	for key, _ := range s.Connect {
		keyList = append(keyList, key)
	}
	keyMap := s.getKeyListByEnv(env, keyList)
	for keyName, oldKeyName := range keyMap {
		allConnect[keyName] = s.Connect[oldKeyName]
	}
	return allConnect
}

func (s *Startup) getAllApiUrl(env string) map[string]*urlStruct {
	allApi := make(map[string]*urlStruct)

	keyList := make([]string, 0)
	for key, _ := range s.Api {
		keyList = append(keyList, key)
	}
	keyMap := s.getKeyListByEnv(env, keyList)
	for keyName, oldKeyName := range keyMap {
		allApi[keyName] = s.Api[oldKeyName]
	}
	return allApi
}

// GetAllApiUrl 取得所有服务的地址
func (s *Startup) GetAllApiUrl(env string) []*urlStruct {
	allApi := s.getAllApiUrl(env)
	urlApiList := make([]*urlStruct, 0)
	for _, one := range allApi {
		urlTemp := new(urlStruct)
		_ = conv.Unmarshal(one, urlTemp)

		{
			urlTemp.PaasName = one.PaasName
			urlTemp.ApiName = one.ApiName
			urlTemp.PaasAPIUrl = one.PaasAPIUrl
			urlTemp.paasAPICfg = new(startupCfg.PaasApiConfig)
			urlTemp.paasAPICfg.Domain = one.paasAPICfg.Domain
			urlTemp.paasAPICfg.Polaris = one.paasAPICfg.Polaris
			urlTemp.paasAPICfg.Auth = one.paasAPICfg.Auth
			apiUrl := one.paasAPICfg.Url(one.ApiName)
			if apiUrl == "" { //如果path设置为空，则表示跳过
				continue
			}
			urlTemp.paasAPICfg.Urls = map[string]string{
				one.ApiName: apiUrl,
			}
		}

		urlApiList = append(urlApiList, urlTemp)
	}
	return urlApiList
}

// getAllCustoms 取得所有自定义的数据
func (s *Startup) getAllCustoms(env string) map[string]interface{} {
	allCustom := make(map[string]interface{})

	keyList := make([]string, 0)
	for key, _ := range s.Custom {
		keyList = append(keyList, key)
	}
	keyMap := s.getKeyListByEnv(env, keyList)
	for keyName, oldKeyName := range keyMap {
		allCustom[keyName] = s.Custom[oldKeyName]
	}

	return allCustom
}

// GetAllConfig 直接yaml配置读取到Gdp的对象里
func (s *Startup) GetAllConfig(env string) *Startup {
	startUp := new(Startup)
	startUp.Mysql = s.getAllMysql(env)
	startUp.Redis = s.getAllRedis(env)
	startUp.Tdmq = s.getAllTdmq(env)

	startUp.Connect = s.getAllConnect(env)
	startUp.Api = s.getAllApiUrl(env)
	startUp.Custom = s.getAllCustoms(env)
	return startUp
}

// GetOneMysql xxx
func (s *Startup) GetOneMysql(env string, keyName string) *startupCfg.MysqlConfig {
	allMysql := s.getAllMysql(env)
	if keyName == "" {
		//如果只有一个，则默认取这个
		if len(allMysql) == 1 {
			for _, one := range allMysql {
				return one
			}
		}
		return nil
	}
	if one, ok := allMysql[keyName]; ok {
		return one
	}
	return nil
}

// GetOneRedis xxx
func (s *Startup) GetOneRedis(env string, keyName string) *startupCfg.RedisConfig {
	allRedis := s.getAllRedis(env)
	if keyName == "" {
		//如果只有一个，则默认取这个
		if len(allRedis) == 1 {
			for _, one := range allRedis {
				return one
			}
		}
		return nil
	}
	if one, ok := allRedis[keyName]; ok {
		return one
	}
	return nil
}

// GetOneTdmq xxx
func (s *Startup) GetOneTdmq(env string, keyName string) *startupCfg.TdmqConfig {
	allTdmq := s.getAllTdmq(env)
	if keyName == "" {
		//如果只有一个，则默认取这个
		if len(allTdmq) == 1 {
			for _, one := range allTdmq {
				return one
			}
		}
		return nil
	}
	if one, ok := allTdmq[keyName]; ok {
		return one
	}

	return nil
}

// GetOneConnect xxx
func (s *Startup) GetOneConnect(env string, keyName string) *conn.Connect {
	if keyName == "" {
		return nil
	}
	allConn := s.getAllConnect(env)
	if one, ok := allConn[keyName]; ok {
		return one
	}
	return nil
}

// GetOneCustom xxx
func (s *Startup) GetOneCustom(env string, keyName string) interface{} {
	if keyName == "" {
		return nil
	}
	allCustom := s.getAllCustoms(env)
	if one, ok := allCustom[keyName]; ok {
		return one
	}
	return nil
}

// GetOneApiUrl xxx
func (s *Startup) GetOneApiUrl(env string, paasName string, apiName string) *urlStruct {
	allApi := s.getAllApiUrl(env)
	mapKey := fmt.Sprintf("%s%s%s", paasName, s.paasSplitString, apiName)
	if one, ok := allApi[mapKey]; ok {
		return one
	}
	return nil
}
