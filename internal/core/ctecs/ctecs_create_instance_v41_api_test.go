package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtecsCreateInstanceV41Api_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtecsCreateInstanceV41Api

	// 构造请求
	var bootDiskIsEncrypt bool = true
	var monitorService bool = true
	var isEncrypt bool = true
	request := &CtecsCreateInstanceV41Request{
		ClientToken:       "4cf2962d-e92c-4c00-9181-cfbb2218636c",
		RegionID:          "bb9fdb42056f11eda1610242ac110002",
		AzName:            "cn-huadong1-jsnj1A-public-ctcloud",
		InstanceName:      "ecm-3300",
		DisplayName:       "ecm-3300",
		FlavorID:          "0824679a-dc86-47dc-a0d3-9c330928f4f6",
		ImageType:         1,
		ImageID:           "9d9e8998-8ed5-43b2-99cb-322f2b8cf6fa",
		BootDiskType:      "SATA",
		BootDiskSize:      40,
		BootDiskIsEncrypt: &bootDiskIsEncrypt,
		BootDiskCmkID:     "3f7e2567-4ed3-4f85-9743-c557d9a94667",
		VpcID:             "4797e8a1-722d-4996-9362-458001813e41",
		OnDemand:          false,
		NetworkCardList: []*CtecsCreateInstanceV41NetworkCardListRequest{
			{
				FixedIP:  "192.168.3.20",
				IsMaster: true,
				SubnetID: "a90eebf0-d798-5017-b9f0-9468bb2301c2",
				NicName:  "net.name",
			},
		},
		ExtIP:        "2",
		ProjectID:    "6732237e53bc4591b0e67d750030ebe3",
		SecGroupList: []string{},
		DataDiskList: []*CtecsCreateInstanceV41DataDiskListRequest{
			{
				DiskMode:  "VBD",
				DiskType:  "SATA",
				DiskSize:  20,
				IsEncrypt: &isEncrypt,
				CmkID:     "3f7e2567-4ed3-4f85-9743-c557d9a94667",
				DiskName:  "ebs.name",
			},
		},
		IpVersion:       "ipv4",
		Bandwidth:       100,
		Ipv6AddressID:   "eip-5sdasd2gfh",
		EipID:           "eip-9jpeyl0frh",
		AffinityGroupID: "259b0c37-1044-41d8-989e",
		KeyPairID:       "c57d0626-8a82-407b-a910-b454907778c3",
		UserPassword:    "1qaz@WSX",
		CycleCount:      6,
		CycleType:       "MONTH",
		AutoRenewStatus: 1,
		UserData:        "ZWNobyBoZWxsbyBnb3N0YWNrIQ==",
		PayVoucherPrice: 20.55,
		LabelList: []*CtecsCreateInstanceV41LabelListRequest{
			{
				LabelKey:   "test-key",
				LabelValue: "test-value",
			},
		},
		GpuDriverKits:       "CUDA 11.4.3 Driver 470.82.01 CUDNN 8.8.1.3",
		MonitorService:      &monitorService,
		InstanceDescription: "云主机描述信息",
		LineType:            "bgp_standalone",
		SecurityProduct:     "EnterpriseEdition",
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
