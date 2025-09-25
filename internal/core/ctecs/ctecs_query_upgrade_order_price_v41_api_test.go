package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtecsQueryUpgradeOrderPriceV41Api_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtecsQueryUpgradeOrderPriceV41Api

	// 构造请求
	request := &CtecsQueryUpgradeOrderPriceV41Request{
		RegionID:        "41f64827f25f468595ffa3a5deb5d15d",
		ResourceUUID:    "bandwidth-czpnl3k1mg",
		ResourceType:    "VM",
		FlavorName:      "s2.medium.2",
		Bandwidth:       101,
		DiskSize:        40,
		NatType:         "large",
		IpPoolBandwidth: 9,
		ElbType:         "standardI",
		CbrValue:        100,
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
