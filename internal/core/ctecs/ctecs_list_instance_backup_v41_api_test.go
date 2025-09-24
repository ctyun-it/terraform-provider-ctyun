package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtecsListInstanceBackupV41Api_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtecsListInstanceBackupV41Api

	// 构造请求
	request := &CtecsListInstanceBackupV41Request{
		RegionID:             "bb9fdb42056f11eda1610242ac110002",
		PageNo:               1,
		PageSize:             10,
		InstanceID:           "de70ef00-1ea0-459a-b74d-b06272561a32",
		RepositoryID:         "de70ef00-1ea0-459a-b74d-b06272561a32",
		InstanceBackupID:     "ed48dc25-d6bb-48e6-b202-3e36ee6321a3",
		QueryContent:         "backup-test01",
		InstanceBackupStatus: "ACTIVE",
		ProjectID:            "0",
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
