package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtvpcUpdateRouteTableRulesAttributeApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtvpcUpdateRouteTableRulesAttributeApi

	// 构造请求
	request := &CtvpcUpdateRouteTableRulesAttributeRequest{
		ClientToken:  "79fa97e3-c48b-xxxx-9f46-6a13d8163678",
		RegionID:     "",
		RouteTableID: "rtb-xxxxx",
		RouteRules: []*CtvpcUpdateRouteTableRulesAttributeRouteRulesRequest{
			{
				NextHopID:   "port-xxxx",
				NextHopType: "havip",
				Destination: "192.168.0.1/32",
				IpVersion:   4,
				Description: "test",
				RouteRuleID: "route-rule-xxx",
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
