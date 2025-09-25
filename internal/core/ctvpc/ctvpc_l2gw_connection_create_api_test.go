package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtvpcL2gwConnectionCreateApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtvpcL2gwConnectionCreateApi

	// 构造请求
	var description string = ""
	request := &CtvpcL2gwConnectionCreateRequest{
		RegionID:    "",
		Name:        "",
		Description: &description,
		L2gwID:      "",
		SubnetID:    "",
		L2conIp:     "",
		TunnelID:    0,
		TunnelIp:    "",
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
