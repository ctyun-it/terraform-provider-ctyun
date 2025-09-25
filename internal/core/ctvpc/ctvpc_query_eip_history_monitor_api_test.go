package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtvpcQueryEipHistoryMonitorApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtvpcQueryEipHistoryMonitorApi

	// 构造请求
	request := &CtvpcQueryEipHistoryMonitorRequest{
		RegionID:    "",
		DeviceIDs:   []string{},
		MetricNames: []string{},
		StartTime:   "",
		EndTime:     "",
		Period:      14400,
		PageNumber:  1,
		PageNo:      1,
		PageSize:    1,
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
