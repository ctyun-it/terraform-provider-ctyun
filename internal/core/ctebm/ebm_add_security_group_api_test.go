package ctebm

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestEbmAddSecurityGroupApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.EbmAddSecurityGroupApi

	// 构造请求
	request := &EbmAddSecurityGroupRequest{
		RegionID:            "81f7728662dd11ec810800155d307d5b",
		AzName:              "az1",
		InstanceUUID:        "ss-eztbrc1j541lb2pim0lxr3uctdfo",
		InterfaceUUID:       "ifg-flu2fonep4dcsxhesj9qsflktfst",
		SecurityGroupIDList: "71386230-c9e7-5465-85fc-29e26f79806d,71386230-c9e7-5465-85fc-29e26f798062",
	}

	// 发起调用
	response, err := api.Do(context.Background(), *credential, request)
	if err != nil {
		t.Log("request error:", err)
		t.Fail()
		return
	}
	t.Logf("%+v\n", *response)
}
