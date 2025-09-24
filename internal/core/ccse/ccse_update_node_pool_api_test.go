package ccse

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCcseUpdateNodePoolApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CcseUpdateNodePoolApi

	// 构造请求
	var enableAutoScale bool = true
	request := &CcseUpdateNodePoolRequest{
		ClusterId:       "47281b02f87757478f20b1827c97cadf",
		NodePoolId:      "3074a5bea28dec198475778f20b1827c",
		RegionId:        "bb9fdb42056f11eda1610242ac110002",
		NodePoolName:    "test",
		BillMode:        "1",
		CycleCount:      1,
		CycleType:       "MONTH",
		AutoRenewStatus: 1,
		Description:     "",
		DataDisks: []*CcseUpdateNodePoolDataDisksRequest{
			{
				Size:         100,
				DiskSpecName: "SAS",
			},
		},
		Labels: &CcseUpdateNodePoolLabelsRequest{},
		Taints: []*CcseUpdateNodePoolTaintsRequest{
			{
				Key:    "",
				Value:  "",
				Effect: "",
			},
		},
		EnableAutoScale:          &enableAutoScale,
		MaxNum:                   9,
		MinNum:                   0,
		SysDiskSize:              0,
		SysDiskType:              "",
		VisibilityPostHostScript: "",
		VisibilityHostScript:     "",
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
