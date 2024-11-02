[TOC]

# Startupconfig

启动配置

## Installation

```bash
go get git.woa.com/gdp/lib/startupconfig
```

## Usage
[examples](api_test.go)

### Basic

```go
import (
    "log"

    startup "git.woa.com/gdp/lib/startupconfig"
)

func main() {
    conf, err := startup.New("config.yaml")
    if err != nil {
        log.Fatal(err)
        return
    }
    // Mysql配置
    mysqlName := "myMysql"
    conf.Mysql(mysqlName) // Database interface
    log.Println(conf.MysqlDriverName(mysqlName))
    log.Println(conf.MysqlDSN(mysqlName)) // root:12345678@tcp(127.0.0.1:3306)/db_database?charset=utf8&parseTime=true&loc=Local
    log.Println(conf.MysqlAddress(mysqlName))
    log.Println(conf.MysqlPassword(mysqlName))
    log.Println(conf.MysqlUser(mysqlName))
    log.Println(conf.MysqlDatabase(mysqlName))
    log.Println(conf.MysqlPassword("noConfig")) // ""
    // Redis配置
    redisName := "myRedis"
    conf.Redis(redisName) // Database interface
    log.Println(conf.RedisAddress(redisName))
    log.Println(conf.RedisDatabase(redisName))
    log.Println(conf.RedisPWD(redisName))
    log.Println(conf.RedisPWD("noConfig")) // "" 
    log.Println(conf.RedisUser(redisName)) // odp 
    log.Println(conf.RedisUseTLS(redisName)) // true 
    // 服务接口配置
    service := "paas1"
    conf.PaasAPI(service)                  // PaasAPI interface
    log.Println(conf.UserPolaris(service)) // 是否可以使用北极星（完整配置了北极星）
    log.Println(conf.PaasAPIDomain(service))
    log.Println(conf.PaasAPIPolaris(service)) // polaris instance
    log.Println(conf.PaasAPIPolarisHost(service))
    log.Println(conf.PaasAPIPolarisNamespace(service))
    log.Println(conf.PaasAPIPolarisService(service))
    log.Println(conf.PaasAPIUrl(service, "orders"))
    log.Println(conf.PaasAPIAuthValue(service, "token"))
    log.Println(conf.PaasAPIAuthValue(service, "tok")) // "",nil
    // TDMQ配置
    tdmqName := "myTDMQ"
    conf.TDMQ(tdmqName) // TDMQ interface
    log.Println(conf.TDMQUrl(tdmqName))
    log.Println(conf.TDMQToken(tdmqName))
    log.Println(conf.TDMQSubscription(tdmqName))
    log.Println(conf.TDMQInitialPosition(tdmqName))
    log.Println(conf.TDMQTopic(tdmqName, "ci"))
    log.Println(conf.TDMQTopic(tdmqName, "cc"))   // ""
    log.Println(conf.TDMQTopic("noConfig", "cc")) // ""
    // 自定义kv配置
    conf.Custom() // Custom config Interface
    log.Println(conf.CustomNormal("timeout"))
    log.Println(conf.CustomSensitive("clientId"))
    log.Println(conf.CustomSensitive("sda")) // "", nil
    // 公共环境变量
    log.Println(conf.ConsulHost())
    log.Println(conf.ConsulToken())
    // Trace配置
    conf.Trace()                    // Trace interface
    log.Println(conf.TraceConfig()) // trace配置实例

    log.Println(conf.StartupConfig())             // 启动配置结构实例
    log.Println(conf.StartupMysqlAll())           // 获取所有mysql配置实例
    log.Println(conf.StartupRedisAll())           // 获取所有redis配置实例
    log.Println(conf.StartupTDMQAll())            // 获取所有tdmq配置实例
    log.Println(conf.StartupCustomNormalAll())    // 获取所有自定义 非加密 配置KV
    log.Println(conf.StartupCustomSensitiveAll()) // 获取所有自定义 加密 配置KV
    log.Println(conf.StartupPaasApiAll())         // 获取所有接口配置实例

    // 将配置内容转换为指定的结构体，加密字段会被自动解密  v1.0.12
    var data map[string]interface{}
    if err := conf.Transform(&data); err != nil {
        log.Fatal(err)
        return
    }
    dataJson, _ := json.Marshal(data)
    log.Println(string(dataJson))
    
    // Supported Since v1.0.17 用户操作记录配置
    conf.Recorder()                 // recorder config instance
    conf.RecorderTDMQ()             // recorder config relate tdmq config instance
    t.Log(conf.RecorderTDMQURL())   // get Recorder TDMQ URL
    t.Log(conf.RecorderTDMQToken()) // get Recorder TDMQ Token
    t.Log(conf.RecorderTDMQTopic()) // get Recorder TDMQ topic, aim topic should be "userOperation"
    t.Log(conf.RecorderGroups())    // get Recorder groups information
    group := conf.RecorderGroup("cd") // get Recorder group of key "cd"
    t.Log(group.GroupId) // get Recorder groupId of key "cd"  
    t.Log(group.Module("pod")) // get Recorder module of key "pod" from group "cd"(key)
}
```

### Advance

**Supported Since v1.1.0**

支持通过 JSON Path 获取配置内容，并转换为自定义的结构

```go
func TestConvert(t *testing.T) {
    conf, err := New("config.yaml")
    if err != nil {
      t.Error(err)
      return
    }
    type Custom struct {
      Max   int               `json:"max"`
      Users map[string]string `json:"users"`
      Key   Decrypted         `json:"key"`
    }
    cs := new(Custom)
    if err := conf.ConvertTo("custom.normal", &cs); err != nil {
      t.Error(err)
      return
    }
    t.Logf("%+v", cs) // &{Max:100 Users:map[a:a b:b c:c] Key:87654321}

    var users map[string]string
    if err := conf.ConvertFromCustomNormalTo("users", &users); err != nil {
      t.Error(err)
      return
    }
    t.Logf("%+v", users) // map[a:a b:b c:c]

    max, err := conf.Decrypted("custom.normal.max")
    if err == nil {
      t.Errorf("unexpected")
      return
    }
    t.Logf("max:%v, err:%s", max, err) // max:, err:encoding/hex: odd length hex string

    key := conf.MustDecrypted("custom.normal.key")
    t.Log(key) // 87654321

    us, err := conf.CustomNormalDecrypted("users")
    if err == nil {
      t.Errorf("unexpected")
      return
    }
    t.Logf("us:%v, err:%s", us, err) // us:, err:encoding/hex: invalid byte: U+007B '{'

    key1 := conf.CustomNormalMustDecrypted("key")
    t.Log(key1) // 87654321

    maxRs := conf.GetValue("custom.normal.max")
    t.Log(maxRs.Int()) // 100

    encrypted := conf.GetValueFromCustomNormal("key")
    t.Log(encrypted.String()) // 30acf6565b803600ceaad1584b477d29a784e8f52644c828dd1be1d0dcabd25b
}
```



##  配置模板

```yaml
api: # 接口调用配置
  paas1: # 接口服务名
    auth: # 加密授权信息
      token: 6d6d4ee0a871d34f47e39450f9219883112839e7b2bfac418f35abef09961317
    domain: http://odp-platform.idle-paas-discovery.odpcldev.woa.com
    polaris:
      host: odp-platform.idle-paas-discovery
      namespace: Development
      service: odpqclouddev120.odp-platform.idle-paas-discovery
    urls:
      orders: /vi/project/{projectName}/paas/{paasName}/orders
      resource: /vi/project/{projectName}/paas/{paasName}/resource
  paas2: # 接口服务名
    auth: # 加密授权信息
      id: 6d6d4ee0a871d34f47e39450f9219883112839e7b2bfac418f35abef09961317
    domain: http://odp-platform.startup-configuration.odpcldev.woa.com
    polaris:
      host: odp-platform.startup-configuration
      namespace: Development
      service: odpqclouddev120.odp-platform.startup-configuration
    urls:
      config: /vi/project/{projectName}/paas/{paasName}/config
mysql: # mysql 配置
  envoyDB:
    address: 127.0.0.1:3306
    database: db_test
    pwEncoded: 6d6d4ee0a871d34f47e39450f9219883112839e7b2bfac418f35abef09961317 # 加密password
    username: root
  myMysql:
    address: 127.0.0.1:3306
    database: db_test
    pwEncoded: 6d6d4ee0a871d34f47e39450f9219883112839e7b2bfac418f35abef09961317 # 加密password
    username: root
redis: # redis 配置
  myRedis:
    address: 127.0.0.1:6666
    database: 0
    pwEncoded: 30acf6565b803600ceaad1584b477d29a784e8f52644c828dd1be1d0dcabd25b # 加密password
    username: odp
    useTLS: true
  otherRedis:
    address: 127.0.0.1:6666
    database: 0
    pwEncoded: 30acf6565b803600ceaad1584b477d29a784e8f52644c828dd1be1d0dcabd25b # 加密password
tdmq: # tdmq 配置
  myTDMQ:
    brokerAddr: http://pulsar-25m59wk2dg85.tdmq-pulsar.ap-sh.qcloud.tencenttdmq.com:5039
    initialPosition: earliest # earliest/lasted
    jwtToken: 30acf6565b803600ceaad1584b477d29a784e8f52644c828dd1be1d0dcabd25b # 加密token
    subscriptionName: odp-platform.idle-paas-discovery
    topics:
      ci: dev/gdp-event-ci
      cd: dev/gdp-event-cd
      userOperation: dev/gdp-event-operations
  otherTDMQ:
    brokerAddr: http://pulsar-25m59wk2dg85.tdmq-pulsar.ap-sh.qcloud.tencenttdmq.com:5039
    initialPosition: earliest # earliest/lasted
    jwtToken: 30acf6565b803600ceaad1584b477d29a784e8f52644c828dd1be1d0dcabd25b # 加密token
    subscriptionName: odp-platform.idle-paas-discovery
    topics:
      ci: dev/gdp-event-ci
      cd: dev/gdp-event-cd
custom: # 自定义配置
  sensitive: # 加密敏感配置
    clientId: 6d6d4ee0a871d34f47e39450f9219883112839e7b2bfac418f35abef09961317
    clientSecret: 30acf6565b803600ceaad1584b477d29a784e8f52644c828dd1be1d0dcabd25b
  normal: # 非加密普通配置
    timeout: 10
    max: 100
    regexp: ^[a-z]([a-z0-9])*$
    users: {"a":"a","b":"b","c":"c"}
tracing: # tracing 配置
  service: "Service-A"
  tenantId: "test"
  address: "test.xxx.woa.com"
  httpPort: "12345"
  grpcPort: "12346"
  sampleRatio: 0.01
userOperation: # 用户操作记录配置
  mqConfigKey: myTDMQ # 对应消息队列配置key，此处对应 tdmq.myTDMQ，对应topic key 必须为 userOperation
  groups:  # 操作分组列表
    cd: # CD 分组
      groupId: CD
      modules: # 操作模块
        container: Container
        pod: Pod
    ci:
      groupId: CI
  resourceNameMap: # 操作资源英文-中文对照表
    scale: 资源规格
    replicas: 副本数
    webshell: webshell
```

##  加密方法

###  接口加密

```http://nstar.datamore.oa.com/cfgcterApi/Encrypt?originstring=TOKEN```

eg.

```http://nstar.datamore.oa.com/cfgcterApi/Encrypt?originstring=87654321```

```{"error_code":200,"error_message":"","result":"30acf6565b803600ceaad1584b477d29a784e8f52644c828dd1be1d0dcabd25b","logNo":"b06c063d7f6f51d414c149dffbab4373"}```

### SDK加密

```go
import (
	"fmt"
	"log"

	startup "git.woa.com/gdp/lib/startupconfig"
)

func main() {
	encrypted, err := startup.Encrypt("1234")
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println(encrypted)
}
```

##  解密方法

### SDK 解密方法

```go
import (
	"fmt"
	"log"

	startup "git.woa.com/gdp/lib/startupconfig"
)

func main() {
	decrypted, err := startup.DecDecrypt("e9c5a70f7aa9127fbfc4245f52942629e45aef272c7e5919021960d47531bb28")
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println(decrypted)
}
```

## 解密字符串类型 Decrypted

*Decrypted*  字符串类型实现了 json、yaml ([*gopkg.in/yaml.v3*](gopkg.in/yaml.v3))的  Unmarshaler 接口(interface) ，**反序列化**时可**自动解密**得到 解密后的字符串。例子如下：

```go
func TestDecrypted(t *testing.T) {
	type people struct {
		Name Decrypted `json:"name" yaml:"name"`
		Age  int       `json:"age" yaml:"age"`
	}
	obj := people{
		Name: "767723c23f0144464f768b6fc292ac3932d2f17300d7657f578c269c65e9ed2d",
		Age:  10,
	}
	js, _ := json.Marshal(obj)

	var data people
	if err := json.Unmarshal(js, &data); err != nil {
		t.Error(err)
		return
	}
	t.Log(data) // {12345678 10}
	if err := yaml.Unmarshal(js, &data); err != nil {
		t.Error(err)
		return
	}
	t.Log(data) // {12345678 10}
}
```

