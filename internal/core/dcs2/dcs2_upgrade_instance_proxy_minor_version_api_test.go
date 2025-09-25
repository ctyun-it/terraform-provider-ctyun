package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestDcs2UpgradeInstanceProxyMinorVersionApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.Dcs2UpgradeInstanceProxyMinorVersionApi

	// 构造请求
	request := &Dcs2UpgradeInstanceProxyMinorVersionRequest{
		RegionId:          "bb9fdb42056f11eda1610242ac110002",
		ProdInstId:        "873894303586321408",
		ProxyMinorVersion: "2.8.4.130",
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
