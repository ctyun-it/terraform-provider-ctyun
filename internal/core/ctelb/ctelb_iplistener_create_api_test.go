package ctelb

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtelbIplistenerCreateApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtelbIplistenerCreateApi

	// 构造请求
	request := &CtelbIplistenerCreateRequest{
		RegionID: "bb9fdb42056f11eda1610242ac110002",
		GwLbID:   "",
		Name:     "acl11",
		Action: &CtelbIplistenerCreateActionRequest{
			RawType: "forward",
		},
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
