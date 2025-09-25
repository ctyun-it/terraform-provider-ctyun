package ctebsbackup

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestEbsbackupListEbsBackupPolicyDisksApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.EbsbackupListEbsBackupPolicyDisksApi

	// 构造请求
	request := &EbsbackupListEbsBackupPolicyDisksRequest{
		RegionID: "81f7728662dd11ec810800155d307d5b",
		PolicyID: "d15e7d402f8f11ed81370242ac110006",
		PageNo:   1,
		PageSize: 10,
		DiskID:   "b9631b82-2087-4e54-b142-5640156ac932",
		DiskName: "ecm",
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
