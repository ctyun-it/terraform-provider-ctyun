package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtecsCreateInstanceBackupV41Api_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtecsCreateInstanceBackupV41Api

	// 构造请求
	request := &CtecsCreateInstanceBackupV41Request{
		RegionID:                  "bb9fdb42056f11eda1610242ac110002",
		InstanceID:                "a628a7d9-ef97-3b16-8a0a-4a794fcdbc39",
		InstanceBackupName:        "api-temp01",
		InstanceBackupDescription: "creat_test01",
		RepositoryID:              "bbf404b5-8990-42a7-a144-c85ff4520276",
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
