package ctelb

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtelbCreateTargetGroupApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtelbCreateTargetGroupApi

	// 构造请求
	request := &CtelbCreateTargetGroupRequest{
		ClientToken:   "79fa97e3-c48b-xxxx-9f46-6a13d8163678",
		Protocol:      "",
		RegionID:      "",
		Name:          "acl11",
		VpcID:         "",
		HealthCheckID: "",
		Algorithm:     "",
		SessionSticky: &CtelbCreateTargetGroupSessionStickyRequest{
			SessionStickyMode: "",
			CookieExpire:      1,
			RewriteCookieName: "test",
			SourceIpTimeout:   1,
		},
		ProxyProtocol: 0,
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
