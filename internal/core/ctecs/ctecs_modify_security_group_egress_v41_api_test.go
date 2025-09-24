package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtecsModifySecurityGroupEgressV41Api_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtecsModifySecurityGroupEgressV41Api

	// 构造请求
	request := &CtecsModifySecurityGroupEgressV41Request{
		RegionID:              "bb9fdb42056f11eda1610242ac110002",
		SecurityGroupID:       "sg-bp67axxxxzb4p",
		SecurityGroupRuleID:   "79fa97e3-c48b-xxxxx-9f46-6a13d8163678",
		Description:           "modify_test",
		ClientToken:           "123e4567-e89b-12d3-a456-426655440000",
		Action:                "accept",
		Priority:              1,
		Protocol:              "ANY",
		RemoteSecurityGroupID: "sg-tolywxbe1f",
		DestCidrIp:            "0.0.0.0/0",
		RemoteType:            0,
		PrefixListID:          "pl_xxxx",
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
