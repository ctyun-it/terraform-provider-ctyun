package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtvpcModifyEndpointIpVersionApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtvpcModifyEndpointIpVersionApi

	// 构造请求
	var ipAddress string = "10.0.0.1"
	var ipv6Address string = "100::7:0:0:647e:5665"
	request := &CtvpcModifyEndpointIpVersionRequest{
		RegionID:    "81f7728662dd11ec810800155d307d5b",
		EndpointID:  "endpoint-srsiebllhc",
		IpVersion:   2,
		IpAddress:   &ipAddress,
		Ipv6Address: &ipv6Address,
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
