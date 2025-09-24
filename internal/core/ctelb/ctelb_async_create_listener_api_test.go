package ctelb

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtelbAsyncCreateListenerApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtelbAsyncCreateListenerApi

	// 构造请求
	var caEnabled bool = false
	var forwardedForEnabled bool = true
	request := &CtelbAsyncCreateListenerRequest{
		ClientToken:         "",
		RegionID:            "",
		LoadBalanceID:       "",
		Name:                "acl11",
		Description:         "",
		Protocol:            "",
		ProtocolPort:        8080,
		CertificateID:       "",
		CaEnabled:           &caEnabled,
		ClientCertificateID: "",
		TargetGroup: &CtelbAsyncCreateListenerTargetGroupRequest{
			Name:      "test",
			Algorithm: "rr",
			Targets: []*CtelbAsyncCreateListenerTargetGroupTargetsRequest{
				{
					InstanceID:   "xxxxxxxxxx",
					ProtocolPort: 80,
					InstanceType: "vm",
					Weight:       1,
					Address:      "192.168.0.1",
				},
			},
			HealthCheck: &CtelbAsyncCreateListenerTargetGroupHealthCheckRequest{
				Protocol:          "",
				Timeout:           2,
				Interval:          5,
				MaxRetry:          2,
				HttpMethod:        "",
				HttpUrlPath:       "/",
				HttpExpectedCodes: "200",
			},
			SessionSticky: &CtelbAsyncCreateListenerTargetGroupSessionStickyRequest{
				SessionType:        "APP_COOKIE",
				CookieName:         "test",
				PersistenceTimeout: 10000,
			},
		},
		AccessControlID:     "",
		AccessControlType:   "",
		ForwardedForEnabled: &forwardedForEnabled,
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
