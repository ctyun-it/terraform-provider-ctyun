package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtvpcUpdateSubnetIPv6StatusApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtvpcUpdateSubnetIPv6StatusApi

	// 构造请求
	var clientToken string = "79fa97e3-c48b-xxxx-9f46-6a13d8163678"
	var projectID string = ""
	request := &CtvpcUpdateSubnetIPv6StatusRequest{
		ClientToken: &clientToken,
		RegionID:    "",
		ProjectID:   &projectID,
		SubnetID:    "vpc-hfw53u96ku",
		EnableIpv6:  false,
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
