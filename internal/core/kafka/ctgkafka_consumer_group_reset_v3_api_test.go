package ctgkafka

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtgkafkaConsumerGroupResetV3Api_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtgkafkaConsumerGroupResetV3Api

	// 构造请求
	request := &CtgkafkaConsumerGroupResetV3Request{
		RegionId:   "bb9fdb42056f11eda1610242ac110002",
		ProdInstId: "68eef42fd8d042bb960d3c3244d9243e",
		GroupName:  "test-group",
		TopicName:  "topic1",
		RawType:    0,
		PartitionShiftList: []*CtgkafkaConsumerGroupResetV3PartitionShiftListRequest{
			{
				Partition: 0,
				ShiftBy:   -10,
			},
		},
		Time: 1571299747516,
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
