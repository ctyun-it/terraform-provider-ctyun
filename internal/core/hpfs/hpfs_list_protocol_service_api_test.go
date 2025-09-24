package hpfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestHpfsListProtocolServiceApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.HpfsListProtocolServiceApi

	// 构造请求
	request := &HpfsListProtocolServiceRequest{
		RegionID:              "81f7728662dd11ec810800155d307d5b",
		AzName:                "az1",
		SfsUID:                "a2d4b927-9511-5a9b-b1de-cfe55779c1ea",
		ProtocolServiceStatus: "available",
		ProtocolSpec:          "general",
		ProtocolType:          "nfs",
		PageSize:              10,
		PageNo:                1,
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
