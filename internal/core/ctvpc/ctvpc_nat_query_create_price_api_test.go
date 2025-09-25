package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtvpcNatQueryCreatePriceApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtvpcNatQueryCreatePriceApi

	// 构造请求
	var description string = "testnat"
	request := &CtvpcNatQueryCreatePriceRequest{
		RegionID:    "79fa97e3-c48b-xxxx-9f46-6a13d8163678",
		VpcID:       "vpc-bp1xxxu",
		Name:        "fortest",
		Spec:        1,
		Description: &description,
		ClientToken: "79fa97e3-c48b-xxxx-9f46-6a13d8163678",
		CycleType:   "month",
		CycleCount:  1,
		AzName:      "az1",
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
