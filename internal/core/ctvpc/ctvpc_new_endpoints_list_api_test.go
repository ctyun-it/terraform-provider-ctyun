package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtvpcNewEndpointsListApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtvpcNewEndpointsListApi

	// 构造请求
	var endpointName string = "test"
	var queryContent string = "tes"
	var endpointID string = "a2009322-41db-466d-bd32-01ff51bff327"
	request := &CtvpcNewEndpointsListRequest{
		RegionID:     "81f7728662dd11ec810800155d307d5b",
		Page:         1,
		PageNo:       1,
		PageSize:     10,
		EndpointName: &endpointName,
		QueryContent: &queryContent,
		EndpointID:   &endpointID,
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
