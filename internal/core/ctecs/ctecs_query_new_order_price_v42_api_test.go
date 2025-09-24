package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtecsQueryNewOrderPriceV42Api_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtecsQueryNewOrderPriceV42Api

	// 构造请求
	request := &CtecsQueryNewOrderPriceV42Request{
		RegionID:     "41f64827f25f468595ffa3a5deb5d15d",
		ResourceType: "VM",
		Count:        1,
		OnDemand:     false,
		CycleType:    "MONTH",
		CycleCount:   6,
		FlavorName:   "s2.small.1",
		ImageUUID:    "7d2922f3-019e-4dbb-ad84-cc8c3497546c",
		SysDiskType:  "SATA",
		SysDiskSize:  50,
		Disks: []*CtecsQueryNewOrderPriceV42DisksRequest{
			{
				DiskType: "SATA",
				DiskSize: 10,
			},
		},
		Bandwidth:       1,
		DiskType:        "SATA",
		DiskSize:        30,
		DiskMode:        "VBD",
		NatType:         "small",
		IpPoolBandwidth: 6,
		DeviceType:      "physical.t4.large",
		AzName:          "az1",
		OrderDisks: []*CtecsQueryNewOrderPriceV42OrderDisksRequest{
			{
				DiskType: "SATA",
				DiskSize: 10,
			},
		},
		ElbType:  "standardI",
		CbrValue: 100,
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
