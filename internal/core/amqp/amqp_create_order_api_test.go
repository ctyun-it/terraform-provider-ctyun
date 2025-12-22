package amqp

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestAmqpCreateOrderApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.AmqpCreateOrderApi

	// 构造请求
	var enableIpv6 bool = false
	request := &AmqpCreateOrderRequest{
		RegionId:        "bb9fdb42056f11eda1610242ac110002",
		ProjectId:       "0",
		ClusterName:     "RabbitMQ-Instance-Test",
		SpecName:        "rabbitmq.2u4g.cluster",
		NodeNum:         3,
		ZoneList:        []string{},
		DiskType:        "FAST-SSD",
		DiskSize:        300,
		VpcId:           "vpc-grqvu4741a",
		SubnetId:        "subnet-gr36jdeyt0",
		SecurityGroupId: "sg-ufrtt04xq1",
		EnableIpv6:      &enableIpv6,
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
