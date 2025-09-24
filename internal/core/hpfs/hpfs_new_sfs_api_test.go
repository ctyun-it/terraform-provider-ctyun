package hpfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestHpfsNewSfsApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.HpfsNewSfsApi

	// 构造请求
	var onDemand bool = true
	request := &HpfsNewSfsRequest{
		ClientToken: "test-new-hpfs-2023-0719-0003-token",
		RegionID:    "81f7728662dd11ec810800155d307d5b",
		ProjectID:   "0",
		SfsType:     "hpfs_perf",
		SfsProtocol: "nfs",
		OnDemand:    &onDemand,
		CycleType:   "",
		CycleCount:  0,
		SfsName:     "hpfs-2023-0711-16092",
		SfsSize:     512,
		AzName:      "az2",
		ClusterName: "nm0001",
		Baseline:    "200",
		Vpc:         "vpc-93gwgagqxd",
		Subnet:      "subnet-g5t9dg9tbr",
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
