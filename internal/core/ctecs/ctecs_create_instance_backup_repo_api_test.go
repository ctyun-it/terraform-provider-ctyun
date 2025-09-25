package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtecsCreateInstanceBackupRepoApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtecsCreateInstanceBackupRepoApi

	// 构造请求
	request := &CtecsCreateInstanceBackupRepoRequest{
		RegionID:        "bb9fdb42056f11eda1610242ac110002",
		ProjectID:       "0",
		RepositoryName:  "repo-test",
		CycleCount:      6,
		CycleType:       "MONTH",
		ClientToken:     "4cf2962d-e92c-4c00-9181-cfbb2218636c",
		Size:            100,
		AutoRenewStatus: 1,
		PayVoucherPrice: 20.55,
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
