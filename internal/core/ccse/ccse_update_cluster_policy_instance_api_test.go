package ccse

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCcseUpdateClusterPolicyInstanceApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CcseUpdateClusterPolicyInstanceApi

	// 构造请求
	request := &CcseUpdateClusterPolicyInstanceRequest{
		ClusterId:        "47281b02f87757478f20b1827c97cadf",
		PolicyName:       "xxx",
		RegionId:         "bb9fdb42056f11eda1610242ac110002",
		InstanceId:       12345,
		PolicyParameters: &CcseUpdateClusterPolicyInstancePolicyParametersRequest{},
		PolicyScope:      "test3",
		PolicyAction:     "deny",
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
