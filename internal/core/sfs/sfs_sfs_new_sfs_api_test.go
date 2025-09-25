package sfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestSfsSfsNewSfsApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.SfsSfsNewSfsApi

	// 构造请求
	var isEncrypt bool = true
	var onDemand bool = true
	request := &SfsSfsNewSfsRequest{
		ClientToken: "参考[请求示例]",
		RegionID:    "参考[请求示例]",
		IsEncrypt:   &isEncrypt,
		KmsUUID:     "参考[请求示例]",
		ProjectID:   "参考[请求示例]",
		SfsType:     "参考[请求示例]",
		SfsProtocol: "参考[请求示例]",
		SfsName:     "参考[请求示例]",
		SfsSize:     500,
		OnDemand:    &onDemand,
		CycleType:   "参考[请求示例]",
		CycleCount:  3,
		AzName:      "参考[请求示例]",
		Vpc:         "参考[请求示例]",
		Subnet:      "参考[请求示例]",
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
