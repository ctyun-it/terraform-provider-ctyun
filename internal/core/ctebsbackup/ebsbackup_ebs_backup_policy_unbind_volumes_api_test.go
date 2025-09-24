package ctebsbackup

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestEbsbackupEbsBackupPolicyUnbindVolumesApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.EbsbackupEbsBackupPolicyUnbindVolumesApi

	// 构造请求
	request := &EbsbackupEbsBackupPolicyUnbindVolumesRequest{
		RegionID:  "81f7728662dd11ec810800155d307d5b",
		PolicyID:  "d15e7d402f8f11ed81370242ac110006",
		VolumeIDs: "9915c3f4-8d78-445a-a1da-d8d9287d506b,9915c3f4-8d78-445a-a1da-d8d9287d506a",
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
