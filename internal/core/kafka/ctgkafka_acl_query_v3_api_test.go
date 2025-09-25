package ctgkafka

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtgkafkaAclQueryV3Api_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtgkafkaAclQueryV3Api

	// 构造请求
	request := &CtgkafkaAclQueryV3Request{
		RegionId:     "bb9fdb42056f11eda1610242ac110002",
		ProdInstId:   "6cd3682a1bdc49bc8ceef0b6f7ea8123",
		Username:     "appUser",
		ResourceType: "TOPIC",
		ResourceName: "topic1",
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
