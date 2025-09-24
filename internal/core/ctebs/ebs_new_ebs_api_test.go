package ctebs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestEbsNewEbsApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.EbsNewEbsApi

	// 构造请求
	var clientToken string = "20230211ebsspec7"
	var multiAttach bool = false
	var isEncrypt bool = false
	var kmsUUID string = "111d979e-5f30-4dd6-a167-c8c8cdd8aa7c"
	var projectID string = "0"
	var onDemand bool = false
	var cycleType string = "month"
	var imageID string = "sjsidfnsdfsf"
	var azName string = "az2"
	var deleteSnapWithEbs bool = false
	var backupID string = "0ae97ef5-6ee2-44af-9d05-1a509b0a1be6"
	request := &EbsNewEbsRequest{
		ClientToken:       &clientToken,
		RegionID:          "81f7728662dd11ec810800155d307d5b",
		MultiAttach:       &multiAttach,
		IsEncrypt:         &isEncrypt,
		KmsUUID:           &kmsUUID,
		ProjectID:         &projectID,
		DiskMode:          "VBD",
		DiskType:          "SATA",
		DiskName:          "ebs-newspec-test0211v7",
		DiskSize:          10,
		OnDemand:          &onDemand,
		CycleType:         &cycleType,
		CycleCount:        1,
		ImageID:           &imageID,
		AzName:            &azName,
		ProvisionedIops:   1,
		DeleteSnapWithEbs: &deleteSnapWithEbs,
		Labels: []*EbsNewEbsLabelsRequest{
			{
				Key:   "32ff",
				Value: "fe33",
			},
		},
		BackupID: &backupID,
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
