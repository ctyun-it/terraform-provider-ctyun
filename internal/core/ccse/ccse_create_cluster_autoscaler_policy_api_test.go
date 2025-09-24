package ccse

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCcseCreateClusterAutoscalerPolicyApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CcseCreateClusterAutoscalerPolicyApi

	// 构造请求
	request := &CcseCreateClusterAutoscalerPolicyRequest{
		ClusterId:           "fa3485c14ca3425fa4dc0d5736a1dc8c",
		RegionId:            "bb9fdb42056f11eda1610242ac110002",
		TextPlainDataString: "",
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
