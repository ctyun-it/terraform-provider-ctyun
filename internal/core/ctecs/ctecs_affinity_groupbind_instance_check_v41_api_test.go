package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtecsAffinityGroupbindInstanceCheckV41Api_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtecsAffinityGroupbindInstanceCheckV41Api

	// 构造请求
	request := &CtecsAffinityGroupbindInstanceCheckV41Request{
		RegionID:        "81f7728662dd11ec810800155d307d5b",
		InstanceID:      "9a1fdc59-b6a2-32bb-780f-3b00653c56f7",
		AffinityGroupID: "1d9de965-3d77-25f1-f521-8d6703280406",
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
