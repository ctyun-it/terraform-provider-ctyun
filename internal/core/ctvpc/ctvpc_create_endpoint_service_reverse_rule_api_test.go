package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtvpcCreateEndpointServiceReverseRuleApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtvpcCreateEndpointServiceReverseRuleApi

	// 构造请求
	request := &CtvpcCreateEndpointServiceReverseRuleRequest{
		ClientToken:       "79fa97e3-c48b-xxxx-9f46-6a13d8163678",
		RegionID:          "81f7728662dd11ec810800155d307d5b",
		EndpointServiceID: "endpser-t3jqijb64e",
		EndpointID:        "6442d4f9-c699-4305-a312-9d6c21d54e36",
		TransitIPAddress:  "192.168.4.12",
		TransitPort:       8080,
		Protocol:          "TCP",
		TargetIPAddress:   "192.168.0.5",
		TargetPort:        9090,
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
