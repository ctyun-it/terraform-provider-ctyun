package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtecsCloudAssistantRunCommandApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtecsCloudAssistantRunCommandApi

	// 构造请求
	var saveCommand bool = true
	var enabledParameter bool = false
	request := &CtecsCloudAssistantRunCommandRequest{
		RegionID:         "88f8888888dd88ec888888888d888d8b",
		InstanceIDs:      "8d8e8888-8ed8-88b8-88cb-888f8b8cf8fa,d8e8888-8ed8-88b8-88cb-888f8b8cf8fb",
		CommandName:      "testName",
		Description:      "testDescription",
		CommandType:      "Shell",
		CommandContent:   "",
		WorkingDirectory: "/home/user",
		Timeout:          60,
		SaveCommand:      &saveCommand,
		EnabledParameter: &enabledParameter,
		DefaultParameter: []*CtecsCloudAssistantRunCommandDefaultParameterRequest{
			{
				Key:         "userid",
				Description: "用户id",
				Value:       "test",
			},
		},
		Parameter: &CtecsCloudAssistantRunCommandParameterRequest{},
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
