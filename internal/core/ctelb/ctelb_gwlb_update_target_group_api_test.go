package ctelb

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtelbGwlbUpdateTargetGroupApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtelbGwlbUpdateTargetGroupApi

	// 构造请求
	request := &CtelbGwlbUpdateTargetGroupRequest{
		RegionID:          "bb9fdb42056f11eda1610242ac110002",
		TargetGroupID:     "tg-xxx",
		Name:              "acl11",
		HealthCheckID:     "hc-xxxx",
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
