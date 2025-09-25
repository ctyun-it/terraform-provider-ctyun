package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtvpcDeleteIPv6GatewayApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtvpcDeleteIPv6GatewayApi

	// 构造请求
	var projectID string = "0"
	request := &CtvpcDeleteIPv6GatewayRequest{
		ClientToken:   "79fa97e3-c48b-xxxx-9f46-6a13d8163678",
		Ipv6GatewayID: "ipv6-csn2porta5",
		ProjectID:     &projectID,
		RegionID:      "xj8g-894g-09oi-po09-12ol-6e6a",
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
