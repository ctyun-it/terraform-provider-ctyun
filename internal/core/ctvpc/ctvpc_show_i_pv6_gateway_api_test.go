package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtvpcShowIPv6GatewayApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtvpcShowIPv6GatewayApi

	// 构造请求
	var projectID string = "0"
	request := &CtvpcShowIPv6GatewayRequest{
		RegionID:      "xj8g-894g-09oi-po09-12ol-6e6a",
		ProjectID:     &projectID,
		Ipv6GatewayID: "igw6-fw2zxq7ug4",
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
