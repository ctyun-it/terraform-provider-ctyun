package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtvpcQueryToCyclePriceApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtvpcQueryToCyclePriceApi

	// 构造请求
	request := &CtvpcQueryToCyclePriceRequest{
		ResourceID:   "eip-xxx,bandwidth-xxx",
		ResourceType: "eip(只支持按宽带计费),bandwidth,elb,private_nat,public_nat",
		RegionID:     "",
		CycleType:    "month,year",
		CycleCount:   1,
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
