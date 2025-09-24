package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtecsQueryRenewOrderPriceV42Api_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtecsQueryRenewOrderPriceV42Api

	// 构造请求
	request := &CtecsQueryRenewOrderPriceV42Request{
		RegionID:     "41f64827f25f468595ffa3a5deb5d15d",
		ResourceType: "VM",
		ResourceID:   "6ff3103f8daf41839d13dafa55e981c1",
		CycleType:    "MONTH",
		CycleCount:   1,
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
