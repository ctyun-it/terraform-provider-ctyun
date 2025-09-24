package ctebsbackup

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestEbsbackupListEbsBackupRepoApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.EbsbackupListEbsBackupRepoApi

	// 构造请求
	var hideExpire bool = true
	var asc bool = true
	request := &EbsbackupListEbsBackupRepoRequest{
		RegionID:       "81f7728662dd11ec810800155d307d5b",
		RepositoryName: "test-repo",
		RepositoryID:   "9915c3f4-8d78-445a-a1da-d8d9287d506b",
		Status:         "active",
		HideExpire:     &hideExpire,
		QueryContent:   "repo",
		PageNo:         1,
		PageSize:       10,
		Asc:            &asc,
		Sort:           "createdDate",
		ProjectID:      "0",
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
