package scaling

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestScalingRuleCreateApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.ScalingRuleCreateApi

	// 构造请求
	request := &ScalingRuleCreateRequest{
		RegionID:      "81f7728662dd11ec810800155d307d5b",
		GroupID:       499,
		Name:          "zjy-policy-test1-4",
		RawType:       1,
		OperateUnit:   1,
		OperateCount:  1,
		Action:        1,
		Cycle:         3,
		Day:           []int32{},
		EffectiveFrom: "",
		EffectiveTill: "",
		ExecutionTime: "",
		Cooldown:      300,
		TriggerObj: &ScalingRuleCreateTriggerObjRequest{
			Name:               "as-alarm-f88f",
			MetricName:         "cpu_util",
			Statistics:         "max",
			ComparisonOperator: "ge",
			Threshold:          50,
			Period:             "5m",
			EvaluationCount:    1,
		},
		TargetObj: &ScalingRuleCreateTargetObjRequest{
			MetricName:              "cpu_util",
			TargetValue:             50,
			ScaleOutEvaluationCount: 3,
			ScaleInEvaluationCount:  10,
			OperateRange:            10,
			DisableScaleIn:          false,
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
