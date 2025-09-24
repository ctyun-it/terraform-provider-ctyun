package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtvpcGwlbUpdateApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtvpcGwlbUpdateApi

	// 构造请求
	var projectID string = ""
	var name string = "acl11"
	request := &CtvpcGwlbUpdateRequest{
		RegionID:  "bb9fdb42056f11eda1610242ac110002",
		ProjectID: &projectID,
		GwLbID:    "",
		Name:      &name,
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
