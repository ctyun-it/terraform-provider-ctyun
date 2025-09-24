package ctebsbackup

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestEbsbackupCreateEbsBackupApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.EbsbackupCreateEbsBackupApi

	// 构造请求
	request := &EbsbackupCreateEbsBackupRequest{
		VolumeID:     "0c582801-6b20-4e3a-956a-f3afbb5e9725",
		RegionID:     "81f7728662dd11ec810800155d307d5b",
		Description:  "test-bak's description",
		RepositoryID: "9915c3f4-8d78-445a-a1da-d8d9287d506b",
		Name:         "test-bak",
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
