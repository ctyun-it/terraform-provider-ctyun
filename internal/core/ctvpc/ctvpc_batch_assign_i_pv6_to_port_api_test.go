package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtvpcBatchAssignIPv6ToPortApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtvpcBatchAssignIPv6ToPortApi

	// 构造请求
	request := &CtvpcBatchAssignIPv6ToPortRequest{
		ClientToken: "79fa97e3-c48b-xxxx-9f46-6a13d8163678",
		RegionID:    "",
		Data: []*CtvpcBatchAssignIPv6ToPortDataRequest{
			{
				NetworkInterfaceID: "port-xx",
				Ipv6AddressesCount: 1,
				Ipv6Addresses:      []*string{},
			},
		},
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
