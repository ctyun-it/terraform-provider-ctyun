package ctebs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestEbsRefundEbsApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.EbsRefundEbsApi

	// 构造请求
	var clientToken string = "refund0211v1"
	var regionID string = "81f7728662dd11ec810800155d307d5b"
	var deleteSnapWithEbs bool = true
	request := &EbsRefundEbsRequest{
		ClientToken:       &clientToken,
		DiskID:            "0ae97ef5-6ee2-44af-9d05-1a509b0a1bxx",
		RegionID:          &regionID,
		DeleteSnapWithEbs: &deleteSnapWithEbs,
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
