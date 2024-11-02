package startupconfig

import (
	"encoding/json"
	"testing"

	"github.com/json-iterator/go"
	"gopkg.in/yaml.v3"
)

func TestNew(t *testing.T) {
	conf, err := New("config.yaml")
	if err != nil {
		t.Error(err)
		return
	}

	mysqlName := "myMysql"
	conf.Mysql(mysqlName) // Database interface
	t.Log(conf.MysqlDriverName(mysqlName))
	t.Log(conf.MysqlDSN(mysqlName)) // root:12345678@tcp(127.0.0.1:3306)/db_database?charset=utf8&parseTime=true&loc=Local
	t.Log(conf.MysqlAddress(mysqlName))
	t.Log(conf.MysqlPassword(mysqlName))
	t.Log(conf.MysqlUser(mysqlName))
	t.Log(conf.MysqlDatabase(mysqlName))
	t.Log(conf.MysqlPassword("test"))

	redisName := "myRedis"
	conf.Redis(redisName) // Database interface
	t.Log(conf.RedisAddress(redisName))
	t.Log(conf.RedisDatabase(redisName))
	t.Log(conf.RedisPWD(redisName))
	t.Log(conf.RedisPWD("test"))
	t.Log(conf.RedisUser(redisName))
	t.Log(conf.RedisUseTLS(redisName))
	t.Log(conf.RedisUseTLS("otherRedis"))

	service := "paas1"
	conf.PaasAPI(service) // PaasAPI interface
	t.Log(conf.UserPolaris(service))
	t.Log(conf.PaasAPIDomain(service))
	t.Log(conf.PaasAPIPolaris(service)) // polaris instance
	t.Log(conf.PaasAPIPolarisHost(service))
	t.Log(conf.PaasAPIPolarisNamespace(service))
	t.Log(conf.PaasAPIPolarisService(service))
	t.Log(conf.PaasAPIUrl(service, "orders"))
	t.Log(conf.PaasAPIAuthValue(service, "token"))
	t.Log(conf.PaasAPIAuthValue(service, "tok"))

	tdmqName := "myTDMQ"
	conf.TDMQ(tdmqName) // TDMQ interface
	t.Log(conf.TDMQUrl(tdmqName))
	t.Log(conf.TDMQToken(tdmqName))
	t.Log(conf.TDMQSubscription(tdmqName))
	t.Log(conf.TDMQInitialPosition(tdmqName))
	t.Log(conf.TDMQTopic(tdmqName, "ci"))
	t.Log(conf.TDMQTopic(tdmqName, "cc"))
	t.Log(conf.TDMQTopic("test", "cc"))
	t.Log(conf.TDMQToken("test"))

	conf.Custom() // Custom config Interface
	t.Log(conf.CustomNormal("timeout"))
	t.Log(conf.CustomNormal("users"))
	t.Log(conf.CustomSensitive("clientId"))
	t.Log(conf.CustomSensitive("sda"))

	//t.Log(conf.ConsulHost())
	//t.Log(conf.ConsulToken())

	t.Log(conf.Trace().TAddress())
	t.Log(conf.TraceConfig())

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	js, _ := json.Marshal(conf.StartupConfig())
	t.Log(string(js))

	t.Log(conf.StartupMysqlAll())           // 获取所有mysql配置实例
	t.Log(conf.StartupRedisAll())           // 获取所有redis配置实例
	t.Log(conf.StartupTDMQAll())            // 获取所有tdmq配置实例
	t.Log(conf.StartupCustomNormalAll())    // 获取所有自定义 非加密 配置KV
	t.Log(conf.StartupCustomSensitiveAll()) // 获取所有自定义 加密 配置KV
	t.Log(conf.StartupPaasApiAll())         // 获取所有接口配置实例

	var data map[string]interface{}
	if err := conf.Transform(&data); err != nil {
		t.Error(err)
		return
	}
	dataJson, _ := json.Marshal(data)
	t.Logf("%s", dataJson)

	// Form v1.0.17
	conf.Recorder()                   // recorder config instance
	conf.RecorderTDMQ()               // recorder config relate tdmq config instance
	t.Log(conf.RecorderTDMQURL())     // get Recorder TDMQ URL
	t.Log(conf.RecorderTDMQToken())   // get Recorder TDMQ Token
	t.Log(conf.RecorderTDMQTopic())   // get Recorder TDMQ topic, aim topic should be "userOperation"
	t.Log(conf.RecorderGroups())      // get Recorder groups information
	group := conf.RecorderGroup("cd") // get Recorder group of key "cd"
	t.Log(group.GroupId)              // get Recorder groupId of key "cd"
	t.Log(group.Module("pod"))        // get Recorder module of key "pod" from group "cd"(key)
}

func TestDecDecrypt(t *testing.T) {
	str, err := DecDecrypt("1762780f68096b10c4e96d4a001d1e56acade74f43c0ab3788e606deea5ee11e")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(str)
}

func TestEncrypt(t *testing.T) {
	str, _ := Encrypt("L1K)bp2WOb4j9w,p-QPXB!LaOSM")
	t.Log(str)
}

func TestEncrypted(t *testing.T) {
	str := Encrypted("767723c23f0144464f768b6fc292ac3932d2f17300d7657f578c269c65e9ed2d")
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	js, _ := json.Marshal(&str)
	t.Log(string(js)) // "12345678"

	yml, err := yaml.Marshal(&str)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(string(yml)) // "12345678"

}

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

	obj1 := people{
		Name: "",
		Age:  10,
	}
	js1, _ := json.Marshal(obj1)

	var data1 people
	if err := json.Unmarshal(js1, &data1); err != nil {
		t.Error(err)
		return
	}
	t.Log(data1) // { 10}
	if err := yaml.Unmarshal(js1, &data1); err != nil {
		t.Error(err)
		return
	}
	t.Log(data1) // { 10}
}

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

func TestNew2(t *testing.T) {
	t.Log(string(key))

}
