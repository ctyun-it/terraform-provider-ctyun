package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtecsCreateInstanceByInstanceBackupApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtecsCreateInstanceByInstanceBackupApi

	// 构造请求
	var monitorService bool = true
	request := &CtecsCreateInstanceByInstanceBackupRequest{
		ClientToken:      "4cf2962d-e92c-4c00-9181-cfbb2218636c",
		RegionID:         "bb9fdb42056f11eda1610242ac110002",
		AzName:           "",
		InstanceName:     "ecm-3300",
		DisplayName:      "ecm-3300",
		InstanceBackupID: "b6e2966d-7b1c-385e-abe4-d940caa273b7",
		FlavorID:         "0824679a-dc86-47dc-a0d3-9c330928f4f6",
		VpcID:            "4797e8a1-722d-4996-9362-458001813e41",
		OnDemand:         false,
		SecGroupList:     []string{},
		NetworkCardList: []*CtecsCreateInstanceByInstanceBackupNetworkCardListRequest{
			{
				NicName:  "nic-0701",
				FixedIP:  "192.168.0.2",
				IsMaster: true,
				SubnetID: "a90eebf0-d798-5017-b9f0-9468bb2301c2",
			},
		},
		ExtIP:           "2",
		IpVersion:       "ipv4",
		Bandwidth:       100,
		Ipv6AddressID:   "ddabc7f0-2121-4121-bf85-ec090b3a73fc",
		EipID:           "eip-9jpeyl0frh",
		AffinityGroupID: "924b95c4-68d9-3fbe-835d-fee46397feda",
		KeyPairID:       "c57d0626-8a82-407b-a910-b454907778c3",
		UserPassword:    "HelloCtyun.13",
		CycleCount:      6,
		CycleType:       "MONTH",
		AutoRenewStatus: 1,
		UserData:        "ZWNobyBoZWxsbyBnb3N0YWNrIQ==",
		ProjectID:       "0",
		PayVoucherPrice: 20.55,
		LabelList: []*CtecsCreateInstanceByInstanceBackupLabelListRequest{
			{
				LabelKey:   "test-key",
				LabelValue: "test-value",
			},
		},
		MonitorService: &monitorService,
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
