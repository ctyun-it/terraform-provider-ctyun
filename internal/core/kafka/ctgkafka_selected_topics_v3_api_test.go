package ctgkafka

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtgkafkaSelectedTopicsV3Api_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtgkafkaSelectedTopicsV3Api

	// 构造请求
	request := &CtgkafkaSelectedTopicsV3Request{
		RegionId:   "bb9fdb42056f11eda1610242ac110002",
		ProdInstId: "68eef42fd8d042bb960d3c3244d9243e",
		GroupName:  "test-group",
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
