package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtvpcCreateSgIngressRuleApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtvpcCreateSgIngressRuleApi

	// 构造请求
	var remoteSecurityGroupID string = ""
	var rawRange string = "8000-9000"
	request := &CtvpcCreateSgIngressRuleRequest{
		RegionID:        "79fa97e3-c48b-xxxx-9f46-6a13d8163678",
		SecurityGroupID: "sg-bp67acxxxazb4p",
		SecurityGroupRules: []*CtvpcCreateSgIngressRuleSecurityGroupRulesRequest{
			{
				Direction:             "ingress",
				RemoteType:            0,
				RemoteSecurityGroupID: &remoteSecurityGroupID,
				Action:                "accept",
				Priority:              100,
				Protocol:              "ANY",
				Ethertype:             "IPv4",
				DestCidrIp:            "0.0.0.0/0",
				RawRange:              &rawRange,
			},
		},
		ClientToken: "79fa97e3-c48b-xxxx-9f46-6a13d8163678",
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
