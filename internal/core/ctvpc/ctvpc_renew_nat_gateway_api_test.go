package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtvpcRenewNatGatewayApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtvpcRenewNatGatewayApi

	// 构造请求
	var payVoucherPrice string = "1"
	request := &CtvpcRenewNatGatewayRequest{
		RegionID:        "5A2CFF0E-5718-xxx5-9D4D-70B3FF3898",
		NatGatewayID:    "nat-dd349df",
		ClientToken:     "79fa97e3-c48b-xxxx-9f46-6a13d8163678",
		CycleType:       "month",
		CycleCount:      1,
		PayVoucherPrice: &payVoucherPrice,
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
