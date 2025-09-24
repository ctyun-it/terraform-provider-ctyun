package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtvpcAddIPv6ToIPv6BandwidthApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtvpcAddIPv6ToIPv6BandwidthApi

	// 构造请求
	request := &CtvpcAddIPv6ToIPv6BandwidthRequest{
		RegionID:    "81f7728662dd11ec810800155d307d5b",
		BandwidthID: "bandwidth-sjwu65pf9t",
		Ip:          "ff06::c3",
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
