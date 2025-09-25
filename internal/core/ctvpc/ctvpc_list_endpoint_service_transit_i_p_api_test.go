package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtvpcListEndpointServiceTransitIPApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtvpcListEndpointServiceTransitIPApi

	// 构造请求
	request := &CtvpcListEndpointServiceTransitIPRequest{
		RegionID:          "81f7728662dd11ec810800155d307d5b",
		Page:              1,
		PageNo:            1,
		PageSize:          10,
		EndpointServiceID: "endpser-t3jqijb64e",
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
