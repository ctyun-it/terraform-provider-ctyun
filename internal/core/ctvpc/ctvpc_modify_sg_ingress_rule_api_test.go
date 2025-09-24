package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtvpcModifySgIngressRuleApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtvpcModifySgIngressRuleApi

	// 构造请求
	var action string = ""
	var protocol string = ""
	var remoteSecurityGroupID string = ""
	var destCidrIp string = ""
	var prefixListID string = ""
	request := &CtvpcModifySgIngressRuleRequest{
		RegionID:              "79fa97e3-c48b-xxxx-9f46-6a13d8163678",
		SecurityGroupID:       "sg-bp67acxxxazb4p",
		SecurityGroupRuleID:   "79fa97e3-c48b-xxxxx-9f46-6a13d8163678",
		ClientToken:           "79fa97e3-c48b-xxxx-9f46-6a13d8163678",
		Action:                &action,
		Priority:              0,
		Protocol:              &protocol,
		RemoteSecurityGroupID: &remoteSecurityGroupID,
		DestCidrIp:            &destCidrIp,
		RemoteType:            0,
		PrefixListID:          &prefixListID,
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
