package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtvpcListEndpointServiceApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtvpcListEndpointServiceApi

	// 构造请求
	var id string = "endpser-9hzaohgug8"
	var endpointServiceName string = "test"
	var queryContent string = "tes"
	request := &CtvpcListEndpointServiceRequest{
		RegionID:            "81f7728662dd11ec810800155d307d5b",
		Page:                1,
		PageNo:              1,
		PageSize:            10,
		Id:                  &id,
		EndpointServiceName: &endpointServiceName,
		QueryContent:        &queryContent,
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
