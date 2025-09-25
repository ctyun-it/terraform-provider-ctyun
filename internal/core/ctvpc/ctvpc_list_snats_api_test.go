package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtvpcListSnatsApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtvpcListSnatsApi

	// 构造请求
	var natGatewayID string = "natgw-1o5sdqb7i2"
	var sNatID string = ""
	var subnetID string = ""
	request := &CtvpcListSnatsRequest{
		RegionID:     "xx73f321ea-62ff-11ec-a8bc-005056898fe0",
		NatGatewayID: &natGatewayID,
		SNatID:       &sNatID,
		SubnetID:     &subnetID,
		PageNumber:   11,
		PageSize:     10,
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
