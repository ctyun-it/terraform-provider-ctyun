package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtvpcIpv6BandwidthQueryRenewPriceApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtvpcIpv6BandwidthQueryRenewPriceApi

	// 构造请求
	request := &CtvpcIpv6BandwidthQueryRenewPriceRequest{
		RegionID:    "5A2CFF0E-5718-xxx5-9D4D-70B3FF3898",
		BandwidthID: "bandwidth-xxxx",
		ClientToken: "79fa97e3-c48b-xxxx-9f46-6a13d8163678",
		CycleType:   "month",
		CycleCount:  1,
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
