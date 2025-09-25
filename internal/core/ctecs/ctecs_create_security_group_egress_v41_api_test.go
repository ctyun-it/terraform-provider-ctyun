package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtecsCreateSecurityGroupEgressV41Api_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtecsCreateSecurityGroupEgressV41Api

	// 构造请求
	request := &CtecsCreateSecurityGroupEgressV41Request{
		RegionID:        "bb9fdb42056f11eda1610242ac110002",
		SecurityGroupID: "sg-bp67axxxxzb4p",
		SecurityGroupRules: []*CtecsCreateSecurityGroupEgressV41SecurityGroupRulesRequest{
			{
				Direction:             "egress",
				RemoteType:            0,
				RemoteSecurityGroupID: "sg-tolywxbe1f",
				Action:                "accept",
				Priority:              100,
				Protocol:              "ANY",
				Ethertype:             "IPv4",
				DestCidrIp:            "0.0.0.0/0",
				Description:           "出方向",
				RawRange:              "8000-9000",
			},
		},
		ClientToken: "123e4567-e89b-12d3-a456-426655440000",
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
