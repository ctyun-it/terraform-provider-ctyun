package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtvpcCreateAclRuleApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtvpcCreateAclRuleApi

	// 构造请求
	var destinationPort string = "20:100"
	var sourcePort string = "20:100"
	request := &CtvpcCreateAclRuleRequest{
		ClientToken: "79fa97e3-c48b-xxxx-9f46-6a13d8163678",
		RegionID:    "81f7728662dd11ec810800155d307d5b",
		AclID:       "acl-r5i4zghgvq",
		Rules: []*CtvpcCreateAclRuleRulesRequest{
			{
				Direction:            "ingress",
				Priority:             100,
				Protocol:             "icmp",
				IpVersion:            "ipv4",
				DestinationPort:      &destinationPort,
				SourcePort:           &sourcePort,
				SourceIpAddress:      "10.1.0.0/24",
				DestinationIpAddress: "10.1.0.0/24",
				Action:               "accept",
				Enabled:              "enable",
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
