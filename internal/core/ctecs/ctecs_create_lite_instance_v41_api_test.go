package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtecsCreateLiteInstanceV41Api_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtecsCreateLiteInstanceV41Api

	// 构造请求
	request := &CtecsCreateLiteInstanceV41Request{
		ClientToken:   "4cf2962d-e92c-4c00-9181-cfbb2218636c",
		RegionID:      "bb9fdb42056f11eda1610242ac110002",
		AzName:        "cn-huadong1-jsnj1A-public-ctcloud",
		DisplayName:   "ecm-3300",
		FlavorSetType: "fix",
		FlavorName:    "lite1.fix.small.1",
		ImageID:       "9d9e8998-8ed5-43b2-99cb-322f2b8cf6fa",
		CycleCount:    6,
		CycleType:     "MONTH",
		IpVersion:     "ipv4",
		BootDiskSize:  40,
		Bandwidth:     5,
		DataDiskList: []*CtecsCreateLiteInstanceV41DataDiskListRequest{
			{
				DiskType: "SATA",
				DiskSize: 50,
			},
		},
		UserPassword:    "1qaz*WSX",
		AutoRenewStatus: 1,
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
