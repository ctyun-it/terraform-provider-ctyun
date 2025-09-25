package scaling

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestScalingRuleUpdateScheduledApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.ScalingRuleUpdateScheduledApi

	// 构造请求
	request := &ScalingRuleUpdateScheduledRequest{
		RegionID:      "81f7728662dd11ec810800155d307d5b",
		GroupID:       489,
		RuleID:        87,
		Name:          "as-policy-ad40",
		OperateUnit:   1,
		OperateCount:  1,
		Action:        2,
		ExecutionTime: "2022-10-13 23:56:00",
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
