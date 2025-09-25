package ctebs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestEbsBatchRollbackEbsSnapApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.EbsBatchRollbackEbsSnapApi

	// 构造请求
	request := &EbsBatchRollbackEbsSnapRequest{
		SnapshotList: "3f868846-f47f-4619-a5b4-a02e9714f744,0fdd318a-89fd-4491-a3f4-6d504bce3c45",
		RegionID:     "81f7728662dd11ec810800155d307d5b",
		AzName:       "az2",
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
