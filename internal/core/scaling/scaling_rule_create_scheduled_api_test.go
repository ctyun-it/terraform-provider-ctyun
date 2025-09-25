package scaling

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestScalingRuleCreateScheduledApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.ScalingRuleCreateScheduledApi

	// 构造请求
	request := &ScalingRuleCreateScheduledRequest{
		RegionID:      "81f7728662dd11ec810800155d",
		GroupID:       471,
		Name:          "as-policy-xudcdcfccd",
		ExecutionTime: "2022-10-17 10:44:00",
		Action:        1,
		OperateUnit:   1,
		OperateCount:  1,
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
