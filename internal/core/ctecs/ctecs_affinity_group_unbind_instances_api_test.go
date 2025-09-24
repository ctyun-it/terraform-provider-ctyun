package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtecsAffinityGroupUnbindInstancesApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtecsAffinityGroupUnbindInstancesApi

	// 构造请求
	request := &CtecsAffinityGroupUnbindInstancesRequest{
		RegionID:        "81f7728662dd11ec810800155d307d5b",
		InstanceIDs:     "d5673536-6c77-8ac8-5b73-19a96fd41dca,191881cf-e766-6909-31a9-c4086b17c1dd",
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
