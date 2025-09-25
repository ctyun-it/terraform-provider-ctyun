package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtvpcQueryCtyunServiceApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtvpcQueryCtyunServiceApi

	// 构造请求
	var vpcID string = "vpc-xxxx"
	var endpointServiceID string = "endpser-9hzaohgug8"
	request := &CtvpcQueryCtyunServiceRequest{
		RegionID:          "81f7728662dd11ec810800155d307d5b",
		PageNumber:        1,
		PageSize:          10,
		VpcID:             &vpcID,
		EndpointServiceID: &endpointServiceID,
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
