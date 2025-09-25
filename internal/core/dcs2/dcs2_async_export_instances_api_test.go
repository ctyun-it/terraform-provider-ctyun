package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestDcs2AsyncExportInstancesApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.Dcs2AsyncExportInstancesApi

	// 构造请求
	request := &Dcs2AsyncExportInstancesRequest{
		RegionId:      "bb9fdb42056f11eda1610242ac110002",
		ProjectId:     "0",
		InstanceName:  "DCS2-blbztb",
		Capacity:      "1",
		ProdInstId:    "610eba1d3dc340dcaa9130b1cd3ccedb",
		Vip:           "192.168.13.83",
		Status:        0,
		EngineVersion: "6.0",
		PayType:       0,
		CpuArchType:   "x86",
		LabelIds:      []string{},
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
