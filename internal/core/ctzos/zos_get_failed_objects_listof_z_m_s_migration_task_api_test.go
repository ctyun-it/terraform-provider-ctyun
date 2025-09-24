package ctzos

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestZosGetFailedObjectsListofZMSMigrationTaskApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.ZosGetFailedObjectsListofZMSMigrationTaskApi

	// 构造请求
	request := &ZosGetFailedObjectsListofZMSMigrationTaskRequest{
		RegionID:    "332232exxxxxxx2e5123866f0",
		MigrationID: "222_mig_a0d60xxxxxxdbeb0213362217328",
		PageSize:    3,
		PageNo:      2,
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
