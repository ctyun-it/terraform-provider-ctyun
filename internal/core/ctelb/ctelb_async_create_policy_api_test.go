package ctelb

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtelbAsyncCreatePolicyApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtelbAsyncCreatePolicyApi

	// 构造请求
	request := &CtelbAsyncCreatePolicyRequest{
		ClientToken: "79fa97e3-c48b-xxxx-9f46-6a13d8163678",
		RegionID:    "",
		ListenerID:  "",
		Name:        "acl11",
		Description: "test",
		Conditions: []*CtelbAsyncCreatePolicyConditionsRequest{
			{
				RuleType:   "PATH",
				MatchType:  "REGEX",
				MatchValue: "/foo",
			},
		},
		TargetGroup: &CtelbAsyncCreatePolicyTargetGroupRequest{
			Name:      "test",
			Algorithm: "rr",
			Targets: []*CtelbAsyncCreatePolicyTargetGroupTargetsRequest{
				{
					InstanceID:   "xxxxxxxxxx",
					ProtocolPort: 80,
					InstanceType: "vm",
					Weight:       1,
					Address:      "192.168.0.1",
				},
			},
			HealthCheck: &CtelbAsyncCreatePolicyTargetGroupHealthCheckRequest{
				Protocol:          "",
				Timeout:           2,
				Interval:          5,
				MaxRetry:          2,
				HttpMethod:        "",
				HttpUrlPath:       "/",
				HttpExpectedCodes: "200",
			},
			SessionSticky: &CtelbAsyncCreatePolicyTargetGroupSessionStickyRequest{
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
