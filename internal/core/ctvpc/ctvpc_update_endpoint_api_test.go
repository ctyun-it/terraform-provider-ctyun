package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtvpcUpdateEndpointApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtvpcUpdateEndpointApi

	// 构造请求
	var endpointName string = "update-it"
	var enableWhitelist bool = false
	var enableDns bool = false
	var whitelist string = "[10.150.0.0/24,10.150.0.0/25]"
	var deleteProtection bool = false
	var protectionService string = ""
	request := &CtvpcUpdateEndpointRequest{
		ClientToken:       "79fa97e3-c48b-xxxx-9f46-6a13d8163678",
		RegionID:          "81f7728662dd11ec810800155d307d5b",
		EndpointID:        "endpoint-tmncs8h97b",
		EndpointName:      &endpointName,
		EnableWhitelist:   &enableWhitelist,
		EnableDns:         &enableDns,
		Whitelist:         &whitelist,
		DeleteProtection:  &deleteProtection,
		ProtectionService: &protectionService,
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
