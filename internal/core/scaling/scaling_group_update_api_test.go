package scaling

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestScalingGroupUpdateApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.ScalingGroupUpdateApi

	// 构造请求
	request := &ScalingGroupUpdateRequest{
		RegionID:            "88f8888888dd88ec888888888d888d8b",
		GroupID:             489,
		Name:                "fcg-group-test",
		MinCount:            1,
		MaxCount:            50,
		ExpectedCount:       1,
		UseLb:               2,
		HealthPeriod:        300,
		SecurityGroupIDList: []string{},
		LbList: []*ScalingGroupUpdateLbListRequest{
			{
				Port:        2235,
				Id:          "lb-la9ik6vb5y",
				Weight:      1,
				HostGroupID: "tg-j36kny3khn",
			},
		},
		SubnetIDList:    []string{},
		MoveOutStrategy: 1,
		RecoveryMode:    1,
		HealthMode:      1,
		ConfigID:        375,
		ConfigList:      []int32{},
		AzStrategy:      1,
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
