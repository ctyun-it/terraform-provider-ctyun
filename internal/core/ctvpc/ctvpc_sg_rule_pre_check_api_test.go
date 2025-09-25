package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtvpcSgRulePreCheckApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtvpcSgRulePreCheckApi

	// 构造请求
	var rawRange string = "8000-9000"
	request := &CtvpcSgRulePreCheckRequest{
		RegionID:        "79fa97e3-c48b-xxxx-9f46-6a13d8163678",
		SecurityGroupID: "sg-bp67acxxxazb4p",
		SecurityGroupRule: &CtvpcSgRulePreCheckSecurityGroupRuleRequest{
			Direction:  "ingress",
			Action:     "accept",
			Priority:   100,
			Protocol:   "ANY",
			Ethertype:  "IPv4",
			DestCidrIp: "0.0.0.0/0",
			RawRange:   &rawRange,
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
