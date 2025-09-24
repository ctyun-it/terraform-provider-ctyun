package ctgkafka

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtgkafkaCreateOrderApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtgkafkaCreateOrderApi

	// 构造请求
	var enableIpv6 bool = false
	var autoRenewStatus bool = false
	request := &CtgkafkaCreateOrderRequest{
		RegionId:            "bb9fdb42056f11eda1610242ac110002",
		ProjectId:           "0",
		CycleCnt:            1,
		ClusterName:         "Kafka-Instance",
		EngineVersion:       "3.6",
		SpecName:            "kafka.8u16g.cluster",
		NodeNum:             3,
		ZoneList:            []string{},
		DiskType:            "FAST-SSD",
		DiskSize:            300,
		VpcId:               "vpc-grqvu4741a",
		SubnetId:            "subnet-gr36jdeyt0",
		SecurityGroupId:     "sg-ufrtt04xq1",
		InstanceNum:         1,
		EnableIpv6:          &enableIpv6,
		PlainPort:           8090,
		SaslPort:            8092,
		SslPort:             8098,
		HttpPort:            8082,
		RetentionHours:      72,
		AutoRenewStatus:     &autoRenewStatus,
		AutoRenewCycleCount: 1,
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
