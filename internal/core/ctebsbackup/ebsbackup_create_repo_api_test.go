package ctebsbackup

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestEbsbackupCreateRepoApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.EbsbackupCreateRepoApi

	// 构造请求
	request := &EbsbackupCreateRepoRequest{
		ClientToken:     "4cf2962d-e92c-4c00-9181-cfbb2218636c",
		RegionID:        "81f7728662dd11ec810800155d307d5b",
		RepositoryName:  "test-repo1",
		Size:            100,
		OnDemand:        "false",
		CycleType:       "MONTH",
		CycleCount:      1,
		AutoRenewStatus: 1,
		ProjectID:       "0",
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
