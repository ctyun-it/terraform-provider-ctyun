package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtvpcQueryResourcesByLabelApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtvpcQueryResourcesByLabelApi

	// 构造请求
	var labelID string = "9b4095099580421ab"
	var labelKey string = "aaaaaa"
	var labelValue string = "aaaaaa"
	request := &CtvpcQueryResourcesByLabelRequest{
		RegionID:   "81f7728662dd11ec810800155d307d5b",
		LabelID:    &labelID,
		LabelKey:   &labelKey,
		LabelValue: &labelValue,
		PageNumber: 0,
		PageSize:   0,
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
