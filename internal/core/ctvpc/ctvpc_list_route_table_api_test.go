package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtvpcListRouteTableApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtvpcListRouteTableApi

	// 构造请求
	var vpcID string = "vpc-xxxx"
	var queryContent string = "xxx"
	var routeTableID string = "rtb-xxxx"
	request := &CtvpcListRouteTableRequest{
		RegionID:     "",
		VpcID:        &vpcID,
		QueryContent: &queryContent,
		RouteTableID: &routeTableID,
		RawType:      0,
		PageNumber:   1,
		PageNo:       1,
		PageSize:     10,
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
