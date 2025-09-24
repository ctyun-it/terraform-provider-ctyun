package ctgkafka

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtgkafkaBindElasticIpApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtgkafkaBindElasticIpApi

	// 构造请求
	request := &CtgkafkaBindElasticIpRequest{
		RegionId:       "bb9fdb42056f11eda1610242ac110002",
		PaasInstanceId: "b9c707b2018c4a5b9ffa5b8c5c837ffc",
		Ip:             "192.168.0.64",
		ElasticIp:      "100.126.8.247",
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
