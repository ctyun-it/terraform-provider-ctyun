package ctgkafka

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtgkafkaUpdateTopicApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtgkafkaUpdateTopicApi

	// 构造请求
	var needFlush bool = false
	var uncleanLeaderElectionEnable bool = false
	var remoteStorageEnable bool = false
	request := &CtgkafkaUpdateTopicRequest{
		RegionId:                    "bb9fdb42056f11eda1610242ac110002",
		ProdInstId:                  "68eef42fd8d042bb960d3c3244d9243e",
		TopicName:                   "test-topic",
		PartitionNum:                8,
		PartitionCapacity:           -1,
		RetentionTime:               3600000,
		MinReplicas:                 2,
		MaxMessage:                  1048576,
		NeedFlush:                   &needFlush,
		TimestampType:               "CreateTime",
		Description:                 "备注",
		StrategyName:                "strategyName",
		CleanupPolicy:               "delete",
		UncleanLeaderElectionEnable: &uncleanLeaderElectionEnable,
		SegmentMs:                   259200000,
		SegmentBytes:                1073741824,
		RemoteStorageEnable:         &remoteStorageEnable,
		LocalRetentionMs:            1073741824,
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
