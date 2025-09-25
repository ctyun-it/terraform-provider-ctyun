package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtvpcL2gwConnectionQueryApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtvpcL2gwConnectionQueryApi

	// 构造请求
	var l2gwID string = "l2gw-hfw53u9"
	var l2ConnectionID string = ""
	request := &CtvpcL2gwConnectionQueryRequest{
		RegionID:       "bb9fdb42056f11eda1610242ac110002",
		L2gwID:         &l2gwID,
		L2ConnectionID: &l2ConnectionID,
		PageNumber:     1,
		PageNo:         1,
		PageSize:       5,
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
