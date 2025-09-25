package ctelb

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtelbUpdateVmPoolAttrApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtelbUpdateVmPoolAttrApi

	// 构造请求
	request := &CtelbUpdateVmPoolAttrRequest{
		RegionID:      "",
		TargetGroupID: "",
		Name:          "acl11",
		HealthCheck: []*CtelbUpdateVmPoolAttrHealthCheckRequest{
			{
				Protocol:          "",
				Timeout:           2,
				Interval:          5,
				MaxRetry:          2,
				HttpMethod:        "",
				HttpUrlPath:       "/",
				HttpExpectedCodes: "200",
			},
		},
		SessionSticky: []*CtelbUpdateVmPoolAttrSessionStickyRequest{
			{
				CookieName:         "test",
				PersistenceTimeout: 10000,
				SessionType:        "APP_COOKIE",
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
