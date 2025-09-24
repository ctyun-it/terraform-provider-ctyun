package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtvpcGwlbUpdateTargetGroupApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtvpcGwlbUpdateTargetGroupApi

	// 构造请求
	var name string = "acl11"
	var healthCheckID string = "hc-xxxx"
	request := &CtvpcGwlbUpdateTargetGroupRequest{
		RegionID:          "bb9fdb42056f11eda1610242ac110002",
		TargetGroupID:     "tg-xxx",
		Name:              &name,
		HealthCheckID:     &healthCheckID,
		SessionStickyMode: 0,
		FailoverType:      0,
		BypassType:        0,
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
