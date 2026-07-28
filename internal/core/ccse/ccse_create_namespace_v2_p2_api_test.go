package ccse

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCcseCreateNamespaceV2P2Api_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CcseCreateNamespaceV2P2Api

	// 构造请求
	request := &CcseCreateNamespaceV2P2Request{
		ClusterName:         "cce-xf6nv8",
		RegionId:            "bb9fdb42056f11eda1610242ac110002",
		TextPlainDataString: "test",
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
