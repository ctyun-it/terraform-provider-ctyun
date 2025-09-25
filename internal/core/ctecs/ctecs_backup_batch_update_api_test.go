package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtecsBackupBatchUpdateApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtecsBackupBatchUpdateApi

	// 构造请求
	request := &CtecsBackupBatchUpdateRequest{
		RegionID: "bb9fdb42056f11eda1610242ac110002",
		UpdateInfo: []*CtecsBackupBatchUpdateUpdateInfoRequest{
			{
				InstanceBackupID:          "b6e2966d-7b1c-385e-abe4-d940caa273b7",
				InstanceBackupName:        "update-test01",
				InstanceBackupDescription: "api_update_test01",
			},
		},
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
