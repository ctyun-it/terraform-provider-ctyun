package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestDcs2DescribePriceApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.Dcs2DescribePriceApi

	// 构造请求
	request := &Dcs2DescribePriceRequest{
		RegionId:      "bb9fdb42056f11eda1610242ac110002",
		ProdInstId:    "b8d6625d63e045e790af8ccfa573d955",
		OrderType:     "INCREASE_SHARDS",
		Version:       "BASIC",
		Edition:       "ClusterSingle",
		ShardMemSize:  "4",
		ShardCount:    3,
		Capacity:      "16",
		CopiesCount:   2,
		HostType:      "S",
		DataDiskType:  "SAS",
		DataDiskSize:  100,
		Period:        "1",
		ChargeType:    "PostPaid",
		EngineVersion: "2.8",
		Size:          1,
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
