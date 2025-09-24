package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtecsCreateVolumeV41Api_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtecsCreateVolumeV41Api

	// 构造请求
	var multiAttach bool = false
	var onDemand bool = false
	var isEncrypt bool = true
	var deleteSnapWithEbs bool = true
	request := &CtecsCreateVolumeV41Request{
		RegionID:          "bb9fdb42056f11eda1610242ac110002",
		DiskMode:          "VBD",
		DiskType:          "SSD",
		DiskName:          "mydisk-000",
		DiskSize:          20,
		ClientToken:       "cbe3840c-bda4-4102-b68f-98c9d7190d69",
		AzName:            "cn-huadong1-jsnj1A-public-ctcloud",
		MultiAttach:       &multiAttach,
		OnDemand:          &onDemand,
		CycleType:         "MONTH",
		CycleCount:        1,
		IsEncrypt:         &isEncrypt,
		KmsUUID:           "3f7e2567-4ed3-4f85-9743-c557d9a94667",
		ProjectID:         "0",
		ImageID:           "8d8e8888-8ed8-88b8-88cb-888f8b8cf8fa",
		ProvisionedIops:   1,
		DeleteSnapWithEbs: &deleteSnapWithEbs,
		Labels: []*CtecsCreateVolumeV41LabelsRequest{
			{
				Key:   "32ff",
				Value: "fe33",
			},
		},
		BackupID: "0ae97ef5-6ee2-44af-9d05-1a509b0a1be6",
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
