package ctebs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestEbsRenewEbsApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.EbsRenewEbsApi

	// 构造请求
	var regionID string = "81f7728662dd11ec810800155d307d5b"
	var clientToken string = "renew-0211v1"
	request := &EbsRenewEbsRequest{
		DiskID:      "0ae97ef5-6ee2-44af-9d05-1a509b0a1bxx",
		RegionID:    &regionID,
		CycleType:   "month",
		CycleCount:  2,
		ClientToken: &clientToken,
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
