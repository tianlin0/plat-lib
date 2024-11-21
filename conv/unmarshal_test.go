package conv

import (
	"encoding/json"
	"fmt"
	cmap "github.com/orcaman/concurrent-map"
	"github.com/tianlin0/plat-lib/cond"
	"github.com/tianlin0/plat-lib/conn"
	jsoniter "github.com/tianlin0/plat-lib/internal/jsoniter/go"
	"log"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestString(t *testing.T) {
	timeStr := "2024-03-05T02:40:19Z"
	var oneTime time.Time
	Unmarshal(timeStr, &oneTime)

	fmt.Println(oneTime)
	fmt.Println(String(oneTime))

	//aa := String(map[string]string{"&": "&&"})
	//fmt.Println(aa)
}

type One struct {
	FullName string
	Persons  []*Two
	Persons2 []Two
	Persons3 []int
	Persons4 []*int
	ClassNum int
	IsTrue   bool `json:"isTrue"`
	OneP     *Two
	TTTP     Two
}
type Two struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestBachGetPageList(t *testing.T) {
	one := new(One)

	aaa := "{\"FullName\":\"aaaaabbbb\"}"

	aaa = "404 no data"

	two := ""

	err := Unmarshal(aaa, one)

	fmt.Println(err, String(one), two)

	return

	var aa int = 3
	Unmarshal(One{
		FullName: "",
		OneP: &Two{
			Name: "bbbb",
			Age:  20,
		},
		Persons: []*Two{
			&Two{
				Name: "bbbb",
				Age:  20,
			},
			&Two{
				Name: "ttttt",
				Age:  40,
			},
		},
		ClassNum: 45,
		Persons3: []int{1, 2, 3},
		Persons4: []*int{&aa},
		Persons2: []Two{
			Two{
				Name: "bbbb",
				Age:  20,
			},
			Two{
				Name: "ttttt",
				Age:  40,
			},
		},
		TTTP: Two{
			Name: "aaaa",
			Age:  33,
		},
		IsTrue: true,
	}, one)
	fmt.Println(String(one))
}

type AuthRoleUser struct {
	RoleId    int64  `json:"role_id" xorm:"not null comment('角色ID') unique(only_one) BIGINT(20)"`
	Userkey   string `json:"userkey" xorm:"not null comment('用户标示，type=0为userid, type=1 为用户组id') unique(only_one) VARCHAR(255)"`
	Type      int    `json:"type" xorm:"comment('可以添加用户组和用户，0为用户，1为用户组') unique(only_one) INT(10)"`
	Domain    string `json:"domain" xorm:"comment('适用的模块，比如项目的角色列表为:project，服务为paas') index(domain) unique(only_one) VARCHAR(255)"`
	Subdomain string `json:"subdomain" xorm:"comment('具体的，比如project为project_id') index(domain) unique(only_one) VARCHAR(255)"`
	IsBool    bool   `json:"isBool" xorm:"comment('具体的，比如project为project_id') index(domain) unique(only_one) VARCHAR(255)"`
}

type AuthRole struct {
	RoleId     int64     `json:"role_id" xorm:"not null pk autoincr BIGINT(20)"`
	RoleName   string    `json:"role_name" xorm:"comment('角色英文名') VARCHAR(255)"`
	RoleCnName string    `json:"role_cn_name" xorm:"comment('角色中文名') VARCHAR(255)"`
	Domain     string    `json:"domain" xorm:"comment('适用的模块，比如项目的角色列表为:project，服务为paas') index(domain) VARCHAR(255)"`
	Subdomain  string    `json:"subdomain" xorm:"comment('具体的子域名，项目为项目ID，服务为具体的服务ID，公共的为空') index(domain) VARCHAR(50)"`
	RoleDesc   string    `json:"role_desc" xorm:"comment('角色说明') VARCHAR(255)"`
	Creator    string    `json:"creator" xorm:"VARCHAR(50)"`
	Createtime time.Time `json:"createtime" xorm:"DATETIME"`
}

type Paas struct {
	PaasId       string    `json:"paas_id" xorm:"not null pk comment('服务ID') VARCHAR(50)"`
	PaasName     string    `json:"paas_name" xorm:"comment('服务英文名，唯一') unique(only_name) VARCHAR(100)"`
	PaasSecret   string    `json:"paas_secret" xorm:"comment('服务密钥，需要加密存储') VARCHAR(1024)"`
	PaasType     string    `json:"paas_type" xorm:"comment('根据类型来判断左侧菜单等') VARCHAR(255)"`
	PaasCnName   string    `json:"paas_cn_name" xorm:"comment('服务中文名') VARCHAR(100)"`
	PaasLanguage string    `json:"paas_language" xorm:"comment('开发语言') VARCHAR(20)"`
	PaasImg      string    `json:"paas_img" xorm:"comment('服务图片') VARCHAR(100)"`
	PaasSummary  string    `json:"paas_summary" xorm:"comment('摘要，简介') LONGTEXT"`
	PaasDesc     string    `json:"paas_desc" xorm:"comment('详细介绍，md信息') LONGTEXT"`
	ProjectId    string    `json:"project_id" xorm:"comment('此服务所属的项目') unique(only_name) VARCHAR(50)"`
	IsRecommend  int       `json:"is_recommend" xorm:"default 0 comment('是否推荐') INT(11)"`
	IsOfficial   int       `json:"is_official" xorm:"default 0 comment('是否官方') INT(11)"`
	IsMarket     int       `json:"is_market" xorm:"default 1 comment('是否显示到应用市场') INT(11)"`
	GitUrl       string    `json:"git_url" xorm:"comment('git全地址') VARCHAR(200)"`
	GitCreator   string    `json:"git_creator" xorm:"default '' comment('git创建时的用户名') VARCHAR(50)"`
	Creator      string    `json:"creator" xorm:"VARCHAR(30)"`
	Createtime   time.Time `json:"createtime" xorm:"DATETIME"`
	Updatetime   time.Time `json:"updatetime" xorm:"DATETIME"`
	IsDelete     int       `json:"is_delete" xorm:"INT(11)"` //1 表示真实删除，2表示有内容删除失败了
	InsertFrom   string    `json:"insert_from" xorm:"default 'odp' comment('创建服务来源，如果是从外部来的话，则自由修改，权限不限制。') VARCHAR(100)"`
}
type PaasInsertParam struct {
	Paas
	IgnoreGit bool   `json:"ignore_git"`
	GitHost   string `json:"git_host"`
}

func TestBachGetPageList2(t *testing.T) {
	//allParam := `{"user_list":[{"userkey":"dsfdsfs","type":0,"isBool":1},{"userkey":"aarenmeng","type":0,"isBool":2}],"role_id":2,"domain":"project","subdomain":"2qkc6haus5c0"}`
	allParam := `{"project_name":"h5ui","paas_cn_name":"h5ui-logserver","ignore_git":0,"paas_desc":"","cluster":"","paas_img":"","git_url":"","paas_name":"h5ui-logserver","env":""}`

	newPaas := PaasInsertParam{}
	_ = Unmarshal(allParam, &newPaas)

	//&{{    h5ui-logserver      0 0 0    0001-01-01 00:00:00 +0000 UTC 0001-01-01 00:00:00 +0000 UTC 0 }    false   }
	//&{{    h5ui-logserver      0 0 0    0001-01-01 00:00:00 +0000 UTC 0001-01-01 00:00:00 +0000 UTC 0 }    false   }

	fmt.Println(newPaas.PaasName, newPaas.PaasCnName)

	//roleUsers := struct {
	//	UserList []*AuthRoleUser `json:"user_list"`
	//}{}
	//_ = Unmarshal(allParam, &roleUsers)
	//
	//fmt.Println(roleUsers)
	//one := new(GdpConfig)
	//Unmarshal(`{"UserProductionAlarmList":["aa","bb",3],"BatchDeleteCdTemplateId":"i31cefk85xe","HostAndPort":{"port":8084,"host":"0.0.0.0"},"ExternalUserContainString":["div","png"]}`, one)
	//fmt.Println(String(one))
}

type CommURL struct {
	GIT_CreateProject string
	Tcr_GetDockerInfo string
}

type GdpConfig struct {
	HostAndPort              *conn.Connect
	APICommUrlMap            *CommURL
	DefaultSystemRoleNameMap map[string][]string
}

func TestRangMap(t *testing.T) {
	newMap := sync.Map{}
	newMap.Store("aaaa", "bbbb")

	tt, err := jsoniter.Marshal(newMap)
	fmt.Println(string(tt), err)

	myMap := cmap.New()

	myMap.Set("key1", "value1")
	myMap.Set("key2", "value2")
	myMap.Set("key3", "value3")

	go func() {
		for i := 1; i < 20; i++ {
			myMap.Set(fmt.Sprintf("%d", i), fmt.Sprintf("value: %d", i))
		}
	}()

	go func() {
		for i := 1; i < 20; i++ {
			fmt.Println(i)
			fmt.Println(String(myMap))
		}
	}()

	select {}

	//Priv = true
	//
	//cfg, err := GetGdpConfig()
	//
	//logs.DefaultLogger().Debug(cfg, err)
}

func TestBachUmailByte111(t *testing.T) {
	var aaa = ""

	mm(func() interface{} {
		aaa = "ok"
		return &aaa
	}, &aaa)

	fmt.Println(aaa)
}

func mm(cc func() interface{}, valuePtr interface{}) {
	retData := cc()

	rf := reflect.ValueOf(valuePtr)
	if rf.Elem().CanSet() {
		fv := reflect.ValueOf(retData)
		if fv.Kind() == reflect.Ptr && fv.Type() == rf.Type() {
			rf.Elem().Set(fv.Elem())
		} else {
			rf.Elem().Set(fv)
		}
	}
}

type CostBizAllInfo struct {
	Desc            string             `json:"desc"`
	Createtime      time.Time          `json:"createtime"`
	Updatetime      time.Time          `json:"updatetime"`
	AllProjectCount int                `json:"all_project_count"`
	OdpProjectList  []*ProjectCostShow `json:"odp_project_list"`
	Creators        *ProjectCostShow   `json:"creators"`
	Managers        []string           `json:"managers"`
}

type CostBizAllInfo2 struct {
	Desc       string          `json:"desc"`
	Createtime time.Time       `json:"createtime"`
	Creators   ProjectCostShow `json:"creators"`
	Managers   []string        `json:"managers"`
}

// ProjectCostShow 业务项目展示
type ProjectCostShow struct {
	Managers   []string  `json:"managers"`
	Createtime time.Time `json:"createtime"`
	Plat       string    `json:"plat"`
}

func TestBachUmailByte333(t *testing.T) {
	//bizList := `[{"ref":"gpid","status":"1","id":"5438","departmentid":"25923","updatetime":"2023-12-02 03:10:00","manager":"zanyzhao","obs_name":"云研发P4","fromplat":"plat","gp_code":"p4","obs_id":"8012","gp_name":"云研发P4","desc":"fdsfdfs","createtime":"2023-12-01 00:10:00","gpid":"301483","department":"技术运营部","odp_project_list":[{"managers":["aaaa","bbbb"],"createtime":"2023-12-02 03:10:00","plat":"mm"}]},{"createtime":"2023-11-30 02:05:00","manager":"lenazhu;sibylswguan","obs_id":"7961","fromplat":"plat","gpid":"301482","obs_name":"代号SE","status":"1","id":"5437","gp_name":"代号SEPC版","gp_code":"DHSEPC","desc":"","ref":"gpid","department":"合作产品部","updatetime":"2023-12-01 05:10:00","departmentid":"57570"},{"gp_name":"代号RA","manager":"jstonesun;stevenjsu;vivianjhe","updatetime":"2023-11-27 23:10:00","ref":"gpid","fromplat":"plat","obs_id":"6132","departmentid":"43448","gpid":"301460","id":"5415","gp_code":"RAMobile","status":"1","desc":"","department":"国内发行线","obs_name":"代号CNC海外","createtime":"2023-11-23 03:05:00"}]`
	//allBizList := make([]*CostBizAllInfo, 0)
	//_ = Unmarshal(bizList, &allBizList)

	//bizList1 := `{"ref":"gpid","status":"1","id":"5438","departmentid":"25923","updatetime":"2023-12-02 03:10:00","manager":"zanyzhao","obs_name":"云研发P4"}`

	//oldBizs := map[string]string{
	//	"desc":       "aaaaa",
	//	"createtime": "2023-11-30 02:05:00",
	//}
	//oldBizs.Desc = "dfdsfdsfdsfds"
	//oldBizs.Createtime = time.Now()
	//allBizs := new(CostBizAllInfo)
	//_ = Unmarshal(oldBizs, allBizs)
	//
	//fmt.Println(String(allBizs))

	//fmt.Println(String(allBizList))

	//bizList2 := `[{"ref":"gpid","status":"1","id":"5438","departmentid":"25923","updatetime":"2023-12-02 03:10:00","manager":"zanyzhao","obs_name":"云研发P4","fromplat":"plat","gp_code":"p4","obs_id":"8012","gp_name":"云研发P4","desc":"fdsfdfs","createtime":"2023-12-01 00:10:00","gpid":"301483","department":"技术运营部","odp_project_list":[{"managers":["aaaa","bbbb"],"createtime":"2023-12-02 03:10:00","plat":"mm"}]},{"createtime":"2023-11-30 02:05:00","manager":"lenazhu;sibylswguan","obs_id":"7961","fromplat":"plat","gpid":"301482","obs_name":"代号SE","status":"1","id":"5437","gp_name":"代号SEPC版","gp_code":"DHSEPC","desc":"","ref":"gpid","department":"合作产品部","updatetime":"2023-12-01 05:10:00","departmentid":"57570"},{"gp_name":"代号RA","manager":"jstonesun;stevenjsu;vivianjhe","updatetime":"2023-11-27 23:10:00","ref":"gpid","fromplat":"plat","obs_id":"6132","departmentid":"43448","gpid":"301460","id":"5415","gp_code":"RAMobile","status":"1","desc":"","department":"国内发行线","obs_name":"代号CNC海外","createtime":"2023-11-23 03:05:00"}]`
	//allBizList2 := make([]CostBizAllInfo, 0)
	//_ = Unmarshal(bizList2, &allBizList2)

	old1 := new(CostBizAllInfo2)
	old1.Creators = ProjectCostShow{
		Plat:       "aaaa",
		Managers:   []string{"aaaa", "bbbb"},
		Createtime: time.Now(),
	}

	allBizList2 := new(CostBizAllInfo)
	_ = Unmarshal(old1, &allBizList2)

	fmt.Println(String(allBizList2))
}

type WorkflowStep struct {
	// Name of the step
	Name string `json:"name,omitempty" protobuf:"bytes,1,opt,name=name"`

	// Template is the name of the template to execute as the step
	Template string `json:"template,omitempty" protobuf:"bytes,2,opt,name=template"`
}

type ParallelSteps struct {
	Steps []WorkflowStep `json:"-" protobuf:"bytes,1,rep,name=steps"`
}

// WorkflowStep is an anonymous list inside of ParallelSteps (i.e. it does not have a key), so it needs its own
// custom Unmarshaller
func (p *ParallelSteps) UnmarshalJSON(value []byte) error {
	// Since we are writing a custom unmarshaller, we have to enforce the "DisallowUnknownFields" requirement manually.

	// First, get a generic representation of the contents
	var candidate []map[string]interface{}
	err := json.Unmarshal(value, &candidate)
	if err != nil {
		return err
	}

	// Generate a list of all the available JSON fields of the WorkflowStep struct
	availableFields := map[string]bool{}
	reflectType := reflect.TypeOf(WorkflowStep{})
	for i := 0; i < reflectType.NumField(); i++ {
		cleanString := strings.ReplaceAll(reflectType.Field(i).Tag.Get("json"), ",omitempty", "")
		availableFields[cleanString] = true
	}

	// Enforce that no unknown fields are present
	for _, step := range candidate {
		for key := range step {
			if _, ok := availableFields[key]; !ok {
				return fmt.Errorf(`json: unknown field "%s"`, key)
			}
		}
	}

	// Finally, attempt to fully unmarshal the struct
	err = json.Unmarshal(value, &p.Steps)
	if err != nil {
		return err
	}
	return nil
}

func (p ParallelSteps) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.Steps)
}

func TestGetCiInfo11(t *testing.T) {
	paramList := new(ParallelSteps)
	//paramList.Steps = []WorkflowStep{
	//	{
	//		Name:     "11111",
	//		Template: "aaaaaa",
	//	}, {
	//		Name:     "2222",
	//		Template: "bbbb",
	//	},
	//}

	jsonStr := `[{"name":"11111"},{"name":"22222"}]`
	err := paramList.UnmarshalJSON([]byte(jsonStr))

	fmt.Println(String(paramList), err)

	//cc := String(paramList)
	//fmt.Println(cc)
	//
	//mm, err := jsoniter.Marshal(paramList)
	//fmt.Println(mm, err)
	//
	//json, err := jsoniter.Marshal(paramList)
	//
	//fmt.Println(json, err)

	//CheckAddrIsUsedInUpstream(context.Background(), "aaaa", "bbbb", "ccccc")

}

type GDPCommonParams struct {
	// ProjectName 项目名
	ProjectName string `json:"projectName" binding:"required"`
	// PaasName 服务名
	PaasName string `json:"paasName" binding:"required"`
	// Operator 操作人
	Operator string `json:"operator" binding:"required"`
}

type GDPLaunchCopyCdWorkflowParams struct {
	GDPCommonParams
	Version   string  `json:"version,omitempty"`
	IsEks     bool    `json:"isEks,omitempty"`
	Gpid      []int64 `json:"gpid,omitempty"`
	SrcCdInfo OneCdInfo
}

type CopyCdCommon struct {
	ProjectName string `json:"project_name"`
	PaasName    string `json:"paasName"`
}

type CopyCdData struct {
	Version int64 `json:"version"`
}

type OneCdInfo struct {
	Cluster string `json:"cluster"`
}

type CopyCd struct {
	CopyCdCommon
	CopyCdData
	SrcCdInfo *OneCdInfo
}

func TestGetCiInfo22(t *testing.T) {

	var params GDPLaunchCopyCdWorkflowParams
	//params.Version = "1"
	params.PaasName = "paasName"
	params.SrcCdInfo.Cluster = "dongbei"

	newCopyCdParam := new(CopyCd)
	_ = Unmarshal(params, newCopyCdParam)

	fmt.Println(String(newCopyCdParam))

}
func TestGetCiInfo33(t *testing.T) {

	var params GDPLaunchCopyCdWorkflowParams
	//params.Version = "1"
	params.PaasName = "paasName123"
	params.SrcCdInfo.Cluster = "dongbei"

	newCopyCdParam := new(CopyCd)
	newCopyCdParam.Version = 55
	_ = Unmarshal(params, newCopyCdParam)

	fmt.Println(String(newCopyCdParam))

}

func getData() (bool, error) {
	return true, nil
}

func TestGetCiInfo44(t *testing.T) {
	aa := true

	mmType := reflect.TypeOf(aa)

	mm := NewPtrByType(mmType)

	fmt.Println(reflect.TypeOf(mm).String())

	params := map[string]interface{}{
		"paasName": "aaabbb",
	}

	newCopyCdParam := new(CopyCd)
	newCopyCdParam.Version = 55

	err := AssignTo(params, newCopyCdParam)
	if err != nil {
		log.Println("Unmarshal error:", err)
	}
	fmt.Println(newCopyCdParam.Version)

	fmt.Println(String(newCopyCdParam))

}

func TestNewPtrByType(t *testing.T) {
	aa := true
	kk := NewPtrByType(reflect.TypeOf(aa))
	fmt.Println(reflect.TypeOf(kk).String())

	bb := 1
	kk = NewPtrByType(reflect.TypeOf(bb))
	fmt.Println(reflect.TypeOf(kk).String())

	cc := "ttt"
	kk = NewPtrByType(reflect.TypeOf(cc))
	fmt.Println(reflect.TypeOf(kk).String())

	dd := map[string]interface{}{}
	kk = NewPtrByType(reflect.TypeOf(dd))
	fmt.Println(reflect.TypeOf(kk).String())

	ee := CopyCd{}
	kk = NewPtrByType(reflect.TypeOf(ee))
	fmt.Println(reflect.TypeOf(kk).String())

	ff := &CopyCd{}
	fff := &ff
	kk = NewPtrByType(reflect.TypeOf(fff))
	fmt.Println(reflect.TypeOf(kk).String())

	gg := fmt.Errorf("mm")
	kk = NewPtrByType(reflect.TypeOf(gg))
	fmt.Println(reflect.TypeOf(kk).String())

	hh := time.Time{}
	kk = NewPtrByType(reflect.TypeOf(hh))
	fmt.Println(reflect.TypeOf(kk).String())

	ii := []byte{}
	kk = NewPtrByType(reflect.TypeOf(ii))
	fmt.Println(reflect.TypeOf(kk).String())

	jj := []int{}
	kk = NewPtrByType(reflect.TypeOf(jj))
	fmt.Println(reflect.TypeOf(kk).String())

	ll := 2.5
	kk = NewPtrByType(reflect.TypeOf(ll))
	fmt.Println(reflect.TypeOf(kk).String())

	mm := []*CopyCd{}
	kk = NewPtrByType(reflect.TypeOf(mm))
	fmt.Println(reflect.TypeOf(kk).String())

}

type newUserAndClientTokenInfo struct {
	XGdpJwtAssertion string `json:"X-Gdp-Jwt-Assertion"`
	AccessToken      string `json:"access_token"`
	ClientID         string `json:"client_id"`
	ClientType       string `json:"client_type"`
	ExpiresIn        int64  `json:"expires_in"`
	RefreshToken     string `json:"refresh_token"`
	TokenType        string `json:"token_type"`
	Userid           string `json:"userid"`
	Usertype         string `json:"usertype"`

	OauthDomain  string
	ClientSecret string
	UserName     string
	UserSecret   string
	ClientScope  string
	CreateTime   time.Time //本机的当前时间
}

func TestDCd(t *testing.T) {
	tokenString := `{"X-Gdp-Jwt-Assertion":"eyJhbGciOiJSUzI1NiIsImtpZCI6ImJlMmEzZTFjLTUwOGEtNDIwNC04Y2I1LWRmYjE3ZmQwZjY1ZSIsInR5cCI6IkpXVCJ9.eyJjbGllbnRfaWQiOiJvZHAtZXh0ZXJuYWwiLCJjbGllbnRfdHlwZSI6ImFkbWluIiwiZXhwIjoxNzAzNzA5NTA1LCJpYXQiOjE3MDM2NjYzMDUsInNjb3BlIjoiZ2RwYWRtaW4ifQ.tRDtK-vEBRGI1cexxybHX5Vp6HFy2hVTCj2cH6K7Gdg63fQzJXyp616WbhEGxo4MU04jBfaNytPiP0HIL2l5Tl5NUHNzXrEsgITP7y-a1ukc2pbKP-nW_MENKcx4i1l8ybGmxDBm-90VcAUZYrcB2IhpkPLNDcFd5pIl3HeBE9pluJUco5fyWPjYN3RVWHfHwmYV--vgWSehsHBlSaPheGvQvG8zwOcm7WJL5Aoz4KeNpIEseSYrw4xvmKZ0F-GpGGdjASO6Cc60qE5RWXkW8m3VPQddN1H1LYVPLDDYhOMstsJhnHIERDkrg8Pxr6VL3lvrsvTxlz4zj70eYjj77A","access_token":"ZJMZYZVLM2YTYWUXYY0ZN2ZILTHLMDATNTUYNGZMOGUWYZUZ","client_id":"odp-external","client_type":"admin","expires_in":600065,"refresh_token":"Y2Q5NZA2ZMYTMTY2MC01ZTU0LTLINZATMGIYNJZLNDBJZWM0","token_type":"Bearer","userid":"odp-external","usertype":"admin","OauthDomain":"","ClientSecret":"827f0a65-48b3-11eb-b993-8e2d46a782b1","UserName":"","UserSecret":"","ClientScope":"odp-external","CreateTime":"2023-12-27 16:38:39"}`

	newClientInfo := new(newUserAndClientTokenInfo)
	err := Unmarshal(tokenString, newClientInfo)

	fmt.Println(err, newClientInfo)

}
func TestTime(t *testing.T) {
	timeNow := time.Time{}

	bb := cond.IsTimeEmpty(timeNow)

	aa := String(timeNow)

	fmt.Println(aa, bb)

}

func TestInt(t *testing.T) {
	intStr := "9223372036854775807"
	aa, err := Int64(intStr)

	bb, err1 := strconv.ParseInt(intStr, 10, 64)

	//aa, err := strconv.ParseInt(intStr, 10, 64)
	fmt.Println(aa, err)
	fmt.Println(bb, err1)
}
func TestInt1(t *testing.T) {
	intFloat := float64(-3.0000)
	aa, err := Int64(intFloat)
	//aa, err := strconv.ParseInt(intStr, 10, 64)
	fmt.Println(intFloat, aa, err)
}
