package scaling

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestScalingRuleUpdateCycleApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.ScalingRuleUpdateCycleApi

	// 构造请求
	request := &ScalingRuleUpdateCycleRequest{
		RegionID:      "81f7728662dd11ec810800155d307d5",
		GroupID:       471,
		RuleID:        899,
		Name:          "as-policy-xcfccd",
		OperateUnit:   1,
		OperateCount:  2,
		Action:        1,
		Cycle:         3,
		Day:           []int32{},
		ExecutionTime: "2022-10-17 10:44:00",
		EffectiveFrom: "2022-10-17 10:44:00",
		EffectiveTill: "2022-11-16 11:44:00",
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
