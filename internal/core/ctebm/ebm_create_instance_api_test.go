package ctebm

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestEbmCreateInstanceApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.EbmCreateInstanceApi

	// 构造请求
	var password string = "****************"
	var systemVolumeRaidUUID string = "r-wtzluqacgzzxgunnabdkpnpjew3d"
	var dataVolumeRaidUUID string = "r-wtzluqacgzzxgunnabdkpnpjew3d"
	var projectID string = "6732237e53bc4591b0e67d750030ebe3"
	var ipType string = "ipv4"
	var bandWidthType string = "standalone"
	var publicIP string = "259b0c37-1044-41d8-989e"
	var securityGroupID string = "259b0c37-1044-41d8-989e-c6f20486c0f4"
	var instanceChargeType string = "ORDER_ON_CYCLE"
	var diskMode string = "VBD"
	var title string = ""
	var title1 string = ""
	var fixedIP string = "192.168.1.1"
	var ipv6 string = ""
	request := &EbmCreateInstanceRequest{
		RegionID:             "100054c0416811e9a6690242ac110002",
		AzName:               "az2",
		DeviceType:           "physical.t3.large",
		Name:                 "pm-3301",
		Hostname:             "host-pm-3301",
		ImageUUID:            "im-xevpi6apqilz1bixmogofyref9qm",
		Password:             &password,
		SystemVolumeRaidUUID: &systemVolumeRaidUUID,
		DataVolumeRaidUUID:   &dataVolumeRaidUUID,
		VpcID:                "4797e8a1-722d-4996-9362-458001813e41",
		ExtIP:                "1",
		ProjectID:            &projectID,
		IpType:               &ipType,
		BandWidth:            100,
		BandWidthType:        &bandWidthType,
		PublicIP:             &publicIP,
		SecurityGroupID:      &securityGroupID,
		DiskList: []*EbmCreateInstanceDiskListRequest{
			{
				DiskType: "data",
				DiskMode: &diskMode,
				Title:    &title,
				RawType:  "SSD",
				Size:     100,
			},
		},
		NetworkCardList: []*EbmCreateInstanceNetworkCardListRequest{
			{
				Title:    &title1,
				FixedIP:  &fixedIP,
				Master:   true,
				Ipv6:     &ipv6,
				SubnetID: "84c95842-13da-47e0-ac94-8fd0861295ad",
			},
		},
		PayVoucherPrice:    20.55,
		AutoRenewStatus:    1,
		InstanceChargeType: &instanceChargeType,
		CycleCount:         6,
		CycleType:          "YEAR",
		OrderCount:         1,
		ClientToken:        "4cf2962d-e92c-4c00-9181-cfbb2218636c",
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
