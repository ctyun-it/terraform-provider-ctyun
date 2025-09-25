package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtvpcUpdateNatGatewayAttributeApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtvpcUpdateNatGatewayAttributeApi

	// 构造请求
	var name string = "acl11"
	var description string = "test"
	request := &CtvpcUpdateNatGatewayAttributeRequest{
		RegionID:     "81f7728662dd11ec810800155d307d5b",
		NatGatewayID: "natgw-1o5sdqb7i2",
		Name:         &name,
		Description:  &description,
		ClientToken:  "79fa97e3-c48b-xxxx-9f46-6a13d8163678",
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
