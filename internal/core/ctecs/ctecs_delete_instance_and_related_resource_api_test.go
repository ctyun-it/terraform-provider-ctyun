package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtecsDeleteInstanceAndRelatedResourceApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtecsDeleteInstanceAndRelatedResourceApi

	// 构造请求
	var deleteVolume bool = true
	var deleteEip bool = true
	request := &CtecsDeleteInstanceAndRelatedResourceRequest{
		RegionID:     "bb9fdb42056f11eda1610242ac110002",
		ClientToken:  "delete-test-001",
		InstanceID:   "755a72c6-ea40-ce04-7ad8-c9f54d38ccfd",
		DeleteVolume: &deleteVolume,
		DeleteEip:    &deleteEip,
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
