package ctelb

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtelbUpdateTargetGroupApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtelbUpdateTargetGroupApi

	// 构造请求
	request := &CtelbUpdateTargetGroupRequest{
		ClientToken:   "79fa97e3-c48b-xxxx-9f46-6a13d8163678",
		RegionID:      "",
		ProjectID:     "",
		ID:            "",
		TargetGroupID: "tg-vzedsj8s49",
		Name:          "acl11",
		HealthCheckID: "",
		Algorithm:     "",
		ProxyProtocol: 0,
		SessionSticky: &CtelbUpdateTargetGroupSessionStickyRequest{
			SessionStickyMode: "",
			CookieExpire:      1,
			RewriteCookieName: "",
			SourceIpTimeout:   10,
		},
		Protocol: "",
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
