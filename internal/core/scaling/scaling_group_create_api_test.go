package scaling

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestScalingGroupCreateApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.ScalingGroupCreateApi

	// 构造请求
	request := &ScalingGroupCreateRequest{
		RegionID:            "81f7728662dd11ec810800155d307d5b",
		SecurityGroupIDList: []string{},
		RecoveryMode:        1,
		Name:                "as-group-901f",
		HealthMode:          1,
		MazInfo: []*ScalingGroupCreateMazInfoRequest{
			{
				MasterId: "subnet-9fr7iesyyp",
				AzName:   "az1",
				OptionId: []string{},
			},
		},
		SubnetIDList:    []string{},
		MoveOutStrategy: 1,
		UseLb:           2,
		VpcID:           "vpc-0wvmt7gh02",
		MinCount:        0,
		MaxCount:        50,
		ExpectedCount:   0,
		HealthPeriod:    300,
		LbList: []*ScalingGroupCreateLbListRequest{
			{
				Port:        2235,
				LbID:        "lb-la9ik6vb5y",
				Weight:      1,
				HostGroupID: "tg-j36kny3khn",
			},
		},
		ProjectID:  "0",
		ConfigID:   375,
		ConfigList: []int32{},
		AzStrategy: 1,
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
