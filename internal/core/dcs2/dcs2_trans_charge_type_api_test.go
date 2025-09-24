package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestDcs2TransChargeTypeApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.Dcs2TransChargeTypeApi

	// 构造请求
	request := &Dcs2TransChargeTypeRequest{
		RegionId:   "bb9fdb42056f11eda1610242ac110002",
		SpuCode:    "DCS2",
		CycleType:  "3",
		CycleCnt:   1,
		ProdInstId: "76d1f92688cc4c22bb907a71debb1c3b",
		AutoPay:    true,
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
