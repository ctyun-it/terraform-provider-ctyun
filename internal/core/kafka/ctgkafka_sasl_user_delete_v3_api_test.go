package ctgkafka

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtgkafkaSaslUserDeleteV3Api_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtgkafkaSaslUserDeleteV3Api

	// 构造请求
	request := &CtgkafkaSaslUserDeleteV3Request{
		RegionId:   "bb9fdb42056f11eda1610242ac110002",
		Username:   "appUser",
		ProdInstId: "00d5f7ee7d9c4f90becb4fe5da5420de",
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
