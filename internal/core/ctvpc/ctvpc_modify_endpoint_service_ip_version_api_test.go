package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtvpcModifyEndpointServiceIpVersionApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtvpcModifyEndpointServiceIpVersionApi

	// 构造请求
	var subnetID string = "interface"
	var ipv6Address string = "100::7:0:0:647e:5665"
	var instanceIDV6 string = "havip-xxxxxx"
	var underlayIp6 string = "100.126.11.95"
	request := &CtvpcModifyEndpointServiceIpVersionRequest{
		RegionID:          "81f7728662dd11ec810800155d307d5b",
		EndpointServiceID: "endpsrv-srsiebllhc",
		IpVersion:         2,
		SubnetID:          &subnetID,
		Ipv6Address:       &ipv6Address,
		InstanceIDV6:      &instanceIDV6,
		UnderlayIp6:       &underlayIp6,
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
