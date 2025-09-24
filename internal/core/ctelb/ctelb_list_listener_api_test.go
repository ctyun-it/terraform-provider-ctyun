package ctelb

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtelbListListenerApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtelbListListenerApi

	// 构造请求
	request := &CtelbListListenerRequest{
		ClientToken:     "",
		RegionID:        "",
		ProjectID:       "",
		IDs:             "'listener-75ex90k9v0,listener-cert-r4cfhgrsss'",
		Name:            "",
		LoadBalancerID:  "",
		AccessControlID: "",
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
