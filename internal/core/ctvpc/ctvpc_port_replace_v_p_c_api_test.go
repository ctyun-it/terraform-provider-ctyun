package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtvpcPortReplaceVPCApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtvpcPortReplaceVPCApi

	// 构造请求
	var ipAddress string = "192.168.0.7"
	request := &CtvpcPortReplaceVPCRequest{
		RegionID:           "81f7728662dd11ec810800155d307d5b",
		NetworkInterfaceID: "port-m8ggj5l97d",
		IpAddress:          &ipAddress,
		SubnetID:           "subnet-a5keks2jdw",
		InstanceID:         "da6358fc-5561-23f5-9bec-1fa912e0aa79",
		SecurityGroupIDs:   []string{},
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
