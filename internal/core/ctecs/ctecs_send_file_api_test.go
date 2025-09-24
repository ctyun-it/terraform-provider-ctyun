package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtecsSendFileApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtecsSendFileApi

	// 构造请求
	var overwrite bool = false
	request := &CtecsSendFileRequest{
		RegionID:        "88f8888888dd88ec888888888d888d8b",
		InstanceIDs:     "8d8e8888-8ed8-88b8-88cb-888f8b8cf8fa,d8e8888-8ed8-88b8-88cb-888f8b8cf8fb",
		FileName:        "testFile.txt",
		Description:     "test file",
		FileContent:     "ZWNobyB0ZXN0",
		TargetDirectory: "/home/user",
		FileOwner:       "root",
		FileGroup:       "root",
		FileMode:        "644",
		Overwrite:       &overwrite,
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
