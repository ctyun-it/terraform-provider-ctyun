package ccse

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCcseDeleteClusterAutoscalerPolicyApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CcseDeleteClusterAutoscalerPolicyApi

	// 构造请求
	request := &CcseDeleteClusterAutoscalerPolicyRequest{
		ClusterId: "47281b02f87757478f20b1827c97cadf",
		Name:      "default-17350973994000003",
		RegionId:  "bb9fdb42056f11eda1610242ac110002",
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
