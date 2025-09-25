package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtecsBatchUpdateInstancesV41Api_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtecsBatchUpdateInstancesV41Api

	// 构造请求
	request := &CtecsBatchUpdateInstancesV41Request{
		RegionID: "88f8888888dd88ec888888888d888d8b",
		UpdateInfo: []*CtecsBatchUpdateInstancesV41UpdateInfoRequest{
			{
				InstanceID:          "8d8e8888-8ed8-88b8-88cb-888f8b8cf8fa",
				DisplayName:         "ecs-0003",
				InstanceName:        "ecm-3300",
				InstanceDescription: "ecm-3300",
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
