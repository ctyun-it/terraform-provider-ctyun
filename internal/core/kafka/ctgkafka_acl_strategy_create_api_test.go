package ctgkafka

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtgkafkaAclStrategyCreateApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtgkafkaAclStrategyCreateApi

	// 构造请求
	request := &CtgkafkaAclStrategyCreateRequest{
		RegionId:   "bb9fdb42056f11eda1610242ac110002",
		ProdInstId: "00d5f7ee7d9c4f90becb4fe5da5420de",
		Name:       "testname",
		Rules: &CtgkafkaAclStrategyCreateRulesRequest{
			Permission: "ALLOW",
			UserName:   "user1",
			Ip:         "192.168.33.21;192.168.30.0/24",
			Operation:  "",
		},
		UseNewTopic: "2",
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
