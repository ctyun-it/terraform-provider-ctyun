package ctelb

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtelbGwlbCreateTargetApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtelbGwlbCreateTargetApi

	// 构造请求
	request := &CtelbGwlbCreateTargetRequest{
		RegionID:      "bb9fdb42056f11eda1610242ac110002",
		TargetGroupID: "",
		InstanceID:    "",
		InstanceType:  "VM",
		Weight:        1,
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
