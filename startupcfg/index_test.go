package startupcfg

import (
	"fmt"
	"github.com/samber/lo"
	"github.com/tianlin0/plat-lib/conv"
	"github.com/tianlin0/plat-lib/encode"
	"strconv"
	"testing"
)

func TestYaml(t *testing.T) {
	config, _ := NewStartupForYamlFile("./config.yaml")
	all := config.GetAllConfig("release")
	fmt.Println(conv.String(all))
	urls := config.GetAllApiUrl("dev")
	fmt.Println(conv.String(urls))

	tdmq := config.GetOneTdmq("release", "TDMQConnect")
	fmt.Println(tdmq)
}

//TofPaasId: odp_notice
//GDPNoticeRobotKey: 02392c8e-cbd0-4dbb-ad5e-ccd309bd4c7b
//GDPCustomerServiceCorpid: wxab249edd27d57738
//GDPCustomerServiceCorpSecret: MGPDWG29k0x-SuuBfwxJ_5aC5CaAuqeEd7Jjf0TTTJM
//GDPGroupSenderRobotKeyList:
//- Key: 02392c8e-cbd0-4dbb-ad5e-ccd309bd4c7b
//Value: gdp-notice
//Label: GDP通知

// TestEncodePass 密码加密
func TestEncodePass(t *testing.T) {

	//dataStr := `{"rrqihrpywawfu":{"PaasId":"rrqihrpywawfu","PaasSecret":"9mCs1wu68PBQ175rPQHXW0PrR9cWOCyB","Version":"","WhiteList":null},"departmentToken":{"PaasId":"departmentToken","PaasSecret":"a73d41d7b206cebc1ee17b17eca19290dcb41107449c4d38971","Version":"","WhiteList":null},"rtx-default":{"PaasId":"b3ede0d25e7d4f0183eb2b413047ec5c","PaasSecret":"27091","Version":"","WhiteList":["odp"]},"rtx-notice":{"PaasId":"a84b16e0d9b247c7945e4a993d6f25b7","PaasSecret":"29053","Version":"","WhiteList":null},"rtx-promethues":{"PaasId":"05822df8b4cb40c5a1582e2869744809","PaasSecret":"29054","Version":"","WhiteList":null},"rtx-monitor":{"PaasId":"92cb1edefa5244b99cb22cd584cfdaf7","PaasSecret":"29055","Version":"","WhiteList":null},"rtx-k8sevent":{"PaasId":"ee8865c6617945a1992683bfa4c45b27","PaasSecret":"29056","Version":"","WhiteList":null},"odp_notice":{"PaasId":"odp_notice","PaasSecret":"pDREbR8GAhZpo7cXTGCMBvw0kP3FNNCf","Version":"","WhiteList":["gdp@tencent.com","odp@tencent.com"]}}`
	//dataBase := base64.StdEncoding.EncodeToString([]byte(dataStr))
	//
	//fmt.Println(dataBase)
	//
	//return

	paasStr := ``
	aaaaa, err := encode.CBCEncrypt(paasStr, "")
	fmt.Println(aaaaa, err)
	bbbb, _ := encode.CBCDecrypt(`dfb6a924c8acebbf3cb5ca8e2fbd62997ffcf7ba12d1f4a8214b47eed4e6d902`, "")
	fmt.Println(bbbb)
}
func TestEncodePass11(t *testing.T) {
	serviceKeys := []int{1, 2, 3}
	eyIDs := lo.Map(serviceKeys, func(serviceKey int, _ int) string {
		return strconv.Itoa(serviceKey)
	})
	fmt.Println(eyIDs)
}

type EncodeType string

func (a EncodeType) Get() string {
	return string(a) + ":result"
}
func TestEncodePass22(t *testing.T) {
	aa := map[string]EncodeType{}

	bb := map[string]string{}
	bb["mm"] = "bb"

	conv.Unmarshal(bb, &aa)

	fmt.Println(aa["mm"].Get())

}

//Ju9FddCBAdTjpHGio7CR
//Ju9FddCBAdTjpHGio7CR
