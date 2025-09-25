package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtvpcCreateHavipApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtvpcCreateHavipApi

	// 构造请求
	var networkID string = "vpc-xxxx"
	var ipAddress string = "192.168.3.1"
	var vipType string = "v4"
	request := &CtvpcCreateHavipRequest{
		ClientToken: "79fa97e3-c48b-xxxx-9f46-6a13d8163678",
		RegionID:    "100054c0416811e9a6690242ac110002",
		NetworkID:   &networkID,
		SubnetID:    "subnet-xxxx",
		IpAddress:   &ipAddress,
		VipType:     &vipType,
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
