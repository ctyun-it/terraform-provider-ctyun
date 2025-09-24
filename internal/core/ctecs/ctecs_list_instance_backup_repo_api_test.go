package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtecsListInstanceBackupRepoApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtecsListInstanceBackupRepoApi

	// 构造请求
	request := &CtecsListInstanceBackupRepoRequest{
		RegionID:       "bb9fdb42056f11eda1610242ac110002",
		ProjectID:      "0",
		RepositoryName: "test-repo",
		RepositoryID:   "9915c3f4-8d78-445a-a1da-d8d9287d506b",
		Status:         "active",
		PageNo:         1,
		PageSize:       5,
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
