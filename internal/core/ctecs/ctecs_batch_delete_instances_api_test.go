package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtecsBatchDeleteInstancesApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtecsBatchDeleteInstancesApi

	// 构造请求
	request := &CtecsBatchDeleteInstancesRequest{
		ClientToken:    "4cf2962d-e92c-4c00-9181-cfbb2218636c",
		RegionID:       "bb9fdb42056f11eda1610242ac110002",
		InstanceIDList: "c7dba5ca-ca72-3823-9429-57a4165600a1,755a72c6-ea40-ce04-7ad8-c9f54d38ccfd",
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
