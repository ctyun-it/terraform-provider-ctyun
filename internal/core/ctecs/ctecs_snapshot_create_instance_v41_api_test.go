package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtecsSnapshotCreateInstanceV41Api_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtecsSnapshotCreateInstanceV41Api

	// 构造请求
	var monitorService bool = true
	request := &CtecsSnapshotCreateInstanceV41Request{
		ClientToken:  "4cf2962d-e92c-4c00-9181-cfbb2218636c",
		RegionID:     "bb9fdb42056f11eda1610242ac110002",
		ProjectID:    "0",
		InstanceName: "ecm-3300",
		DisplayName:  "ecm-3300",
		SnapshotID:   "f02bef68-7aba-f51d-6ead-7d47dde57b45",
		VpcID:        "4797e8a1-722d-4996-9362-458001813e41",
		OnDemand:     false,
		SecGroupList: []string{},
		NetworkCardList: []*CtecsSnapshotCreateInstanceV41NetworkCardListRequest{
			{
				NicName:  "",
				FixedIP:  "",
				IsMaster: true,
				SubnetID: "a90eebf0-d798-5017-b9f0-9468bb2301c2",
			},
		},
		ExtIP:           "2",
		IpVersion:       "ipv4",
		Bandwidth:       100,
		Ipv6AddressID:   "",
		EipID:           "eip-9jpeyl0frh",
		AffinityGroupID: "259b0c37-1044-41d8-989e",
		KeyPairID:       "c57d0626-8a82-407b-a910-b454907778c3",
		UserPassword:    "",
		CycleCount:      6,
		CycleType:       "MONTH",
		AutoRenewStatus: 1,
		UserData:        "ZWNobyBoZWxsbyBnb3N0YWNrIQ==",
		LabelList: []*CtecsSnapshotCreateInstanceV41LabelListRequest{
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
