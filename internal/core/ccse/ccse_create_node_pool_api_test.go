package ccse

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCcseCreateNodePoolApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CcseCreateNodePoolApi

	// 构造请求
	var useAffinityGroup bool = false
	var syncNodeLabels bool = false
	var syncNodeTaints bool = false
	var nodeUnschedulable bool = false
	var enableAutoScale bool = true
	request := &CcseCreateNodePoolRequest{
		ClusterId:                "47281b02f87757478f20b1827c97cadf",
		RegionId:                 "bb9fdb42056f11eda1610242ac110002",
		NodePoolName:             "",
		Description:              "",
		BillMode:                 "1",
		CycleCount:               0,
		CycleType:                "MONTH",
		AutoRenewStatus:          0,
		VisibilityPostHostScript: "",
		VisibilityHostScript:     "",
		DefinedHostnameEnable:    1,
		HostNamePrefix:           "",
		HostNamePostfix:          "",
		ImageType:                1,
		ImageName:                "Centos-79-CCND-CCSE-40-07-amd",
		ImageUuid:                "a8f857a8-ee19-4dc8-a193-819e8ed10dc9",
		LoginType:                "password",
		EcsPasswd:                "",
		KeyName:                  "KeyPair-a589",
		KeyPairId:                "ba425a97-9ad9-2d45-e21a-770d2ebeb477",
		UseAffinityGroup:         &useAffinityGroup,
		AffinityGroupUuid:        "",
		ResourceLabels:           &CcseCreateNodePoolResourceLabelsRequest{},
		SyncNodeLabels:           &syncNodeLabels,
		SyncNodeTaints:           &syncNodeTaints,
		NodeUnschedulable:        &nodeUnschedulable,
		Labels:                   &CcseCreateNodePoolLabelsRequest{},
		Taints: []*CcseCreateNodePoolTaintsRequest{
			{
				Key:    "aaa",
				Value:  "bbb",
				Effect: "NoSchedule",
			},
		},
		VmSpecName:      "c7.xlarge.2",
		VmType:          "C7",
		Cpu:             4,
		Memory:          8,
		MaxNum:          0,
		MinNum:          0,
		EnableAutoScale: &enableAutoScale,
		DataDisks: []*CcseCreateNodePoolDataDisksRequest{
			{
				Size:         100,
				DiskSpecName: "SATA",
			},
		},
		MaxPodNum:   0,
		Gpu:         0,
		SysDiskType: "SATA",
		SysDiskSize: 0,
		AzInfo: []*CcseCreateNodePoolAzInfoRequest{
			{
				AzName: "cn-xinan1-1A",
			},
		},
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
