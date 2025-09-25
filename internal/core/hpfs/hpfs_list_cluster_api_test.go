package hpfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestHpfsListClusterApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.HpfsListClusterApi

	// 构造请求
	request := &HpfsListClusterRequest{
		RegionID:      "81f7728662dd11ec810800155d307d5b",
		SfsType:       "hpfs_perf",
		AzName:        "az2",
		EbmDeviceType: "physical.lcas910b.2xlarge1",
		PageNo:        1,
		PageSize:      10,
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
