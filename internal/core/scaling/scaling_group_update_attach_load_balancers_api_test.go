package scaling

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestScalingGroupUpdateAttachLoadBalancersApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.ScalingGroupUpdateAttachLoadBalancersApi

	// 构造请求
	request := &ScalingGroupUpdateAttachLoadBalancersRequest{
		RegionID: "81f7728662dd11ec810800155d307d5b",
		GroupID:  483,
		LbList: []*ScalingGroupUpdateAttachLoadBalancersLbListRequest{
			{
				Port:        2235,
				LbID:        "lb-la9ik6vb5y",
				Weight:      1,
				HostGroupID: "tg-j36kny3khn",
			},
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
