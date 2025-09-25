package sfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestSfsSfsNewApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.SfsSfsNewApi

	// 构造请求
	var isEncrypt bool = false
	var onDemand bool = false
	request := &SfsSfsNewRequest{
		ClientToken: "",
		RegionID:    "100054c0416811e9a6690242ac110002",
		IsEncrypt:   &isEncrypt,
		KmsUUID:     "",
		ProjectID:   "",
		SfsType:     "",
		SfsProtocol: "",
		Name:        "",
		SfsSize:     0,
		OnDemand:    &onDemand,
		CycleType:   "",
		CycleCount:  0,
		AzName:      "",
		Vpc:         "",
		Subnet:      "",
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
