package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestDcs2UploadSyncRunningLogApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.Dcs2UploadSyncRunningLogApi

	// 构造请求
	request := &Dcs2UploadSyncRunningLogRequest{
		RegionId:   "bb9fdb42056f11eda1610242ac110002",
		TaskId:     "1d22b4aeacf54d4f9096ea735cadc19c",
		SearchDate: "2024-12-01",
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
