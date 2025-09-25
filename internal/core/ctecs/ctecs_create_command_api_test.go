package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtecsCreateCommandApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtecsCreateCommandApi

	// 构造请求
	var enabledParameter bool = false
	request := &CtecsCreateCommandRequest{
		RegionID:         "88f8888888dd88ec888888888d888d8b",
		CommandName:      "testName",
		Description:      "testDescription",
		CommandType:      "Shell",
		CommandContent:   "ZWNobyB0ZXN0",
		WorkingDirectory: "/home/user",
		Timeout:          60,
		EnabledParameter: &enabledParameter,
		DefaultParameter: []*CtecsCreateCommandDefaultParameterRequest{
			{
				Key:         "userid",
				Description: "用户id",
				Value:       "test",
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
