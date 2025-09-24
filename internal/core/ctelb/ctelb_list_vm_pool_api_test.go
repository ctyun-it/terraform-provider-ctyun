package ctelb

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtelbListVmPoolApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtelbListVmPoolApi

	// 构造请求
	request := &CtelbListVmPoolRequest{
		RegionID:      "81f7728662dd11ec810800155d307d5b",
		TargetGroupID: "edbe1a3e-bac3-48a5-9357-d9239d7d1577",
		Name:          "yacos_test_target_group3",
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
