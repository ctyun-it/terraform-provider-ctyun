package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtvpcCreateIPv6BandwidthApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtvpcCreateIPv6BandwidthApi

	// 构造请求
	var payVoucherPrice string = "1"
	request := &CtvpcCreateIPv6BandwidthRequest{
		ClientToken:     "79fa97e3-c48b-xxxx-9f46-6a13d8163678",
		RegionID:        "81f7728662dd11ec810800155d307d5b",
		Bandwidth:       1,
		CycleType:       "year",
		CycleCount:      1,
		Name:            "acl11",
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
