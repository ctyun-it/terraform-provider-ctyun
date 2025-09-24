package scaling

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestScalingRuleUpdateApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.ScalingRuleUpdateApi

	// 构造请求
	var disableScaleIn bool = false
	request := &ScalingRuleUpdateRequest{
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
		Cooldown:      300,
		TriggerObj: &ScalingRuleUpdateTriggerObjRequest{
			Name:               "as-alarm-f8d8fdf",
			MetricName:         "cpu_util",
			Statistics:         "max",
			ComparisonOperator: "ge",
			Threshold:          50,
			Period:             "5m",
			EvaluationCount:    1,
		},
		TargetObj: &ScalingRuleUpdateTargetObjRequest{
			MetricName:              "cpu_util",
			TargetValue:             50,
			ScaleOutEvaluationCount: 3,
			ScaleInEvaluationCount:  10,
			OperateRange:            10,
			DisableScaleIn:          &disableScaleIn,
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
