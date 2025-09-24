package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtecsUpdateEcsLabelV41Api_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtecsUpdateEcsLabelV41Api

	// 构造请求
	request := &CtecsUpdateEcsLabelV41Request{
		RegionID:   "81f7728662dd11ec810800155d307d5b",
		InstanceID: "7ddde90c-1bef-0661-2480-30c536427d09",
		Action:     "ADD",
		LabelList: []*CtecsUpdateEcsLabelV41LabelListRequest{
			{
				LabelKey:   "ctyun",
				LabelValue: "hello",
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
