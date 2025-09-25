package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestDcs2SwitchInstanceHAApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.Dcs2SwitchInstanceHAApi

	// 构造请求
	request := &Dcs2SwitchInstanceHARequest{
		RegionId:     "bb9fdb42056f11eda1610242ac110002",
		ProdInstId:   "873209915236156416",
		MasterNodeIp: "10.50.208.7:25324",
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
