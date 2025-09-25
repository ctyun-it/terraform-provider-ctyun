package ctelb

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtelbCreateListenerApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtelbCreateListenerApi

	// 构造请求
	var caEnabled bool = false
	var forwardedForEnabled bool = true
	request := &CtelbCreateListenerRequest{
		ClientToken:         "79fa97e3-c48b-xxxx-9f46-6a13d8163678",
		RegionID:            "",
		LoadBalancerID:      "",
		Name:                "acl11",
		Description:         "",
		Protocol:            "",
		ProtocolPort:        8080,
		CertificateID:       "",
		CaEnabled:           &caEnabled,
		ClientCertificateID: "",
		DefaultAction: &CtelbCreateListenerDefaultActionRequest{
			RawType: "",
			ForwardConfig: &CtelbCreateListenerDefaultActionForwardConfigRequest{
				TargetGroups: []*CtelbCreateListenerDefaultActionForwardConfigTargetGroupsRequest{
					{
						TargetGroupID: "",
						Weight:        100,
					},
				},
			},
			RedirectListenerID: "",
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
