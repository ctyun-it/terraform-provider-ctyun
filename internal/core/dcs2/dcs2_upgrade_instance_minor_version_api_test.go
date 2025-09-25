package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestDcs2UpgradeInstanceMinorVersionApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.Dcs2UpgradeInstanceMinorVersionApi

	// 构造请求
	request := &Dcs2UpgradeInstanceMinorVersionRequest{
		RegionId:           "bb9fdb42056f11eda1610242ac110002",
		ProdInstId:         "b5fcacfc2e7069553759558b9a4eb27a",
		MinorEngineVersion: "6.2.12.4",
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
