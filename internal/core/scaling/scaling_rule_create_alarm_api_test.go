package scaling

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestScalingRuleCreateAlarmApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.ScalingRuleCreateAlarmApi

	// 构造请求
	request := &ScalingRuleCreateAlarmRequest{
		RegionID:     "81f7728662dd11ec810800155d307d5b",
		GroupID:      511,
		Name:         "as-policy-xudcdcfccd",
		OperateUnit:  1,
		OperateCount: 1,
		Action:       1,
		Cooldown:     300,
		TriggerObj: &ScalingRuleCreateAlarmTriggerObjRequest{
			Name:               "as-alarm-f8d8fdf",
			MetricName:         "cpu_util",
			Statistics:         "max",
			ComparisonOperator: "ge",
			Threshold:          50,
			Period:             "5m",
			EvaluationCount:    1,
		},
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
