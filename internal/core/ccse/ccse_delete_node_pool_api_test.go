package ccse

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCcseDeleteNodePoolApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CcseDeleteNodePoolApi

	// 构造请求
	request := &CcseDeleteNodePoolRequest{
		ClusterId:  "47281b02f87757478f20b1827c97cadf",
		NodePoolId: "3074a5bea28dec198475778f20b1827c",
		RegionId:   "bb9fdb42056f11eda1610242ac110002",
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
