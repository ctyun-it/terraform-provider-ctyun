package ctimage

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtimageCopyImageApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtimageCopyImageApi

	// 构造请求
	request := &CtimageCopyImageRequest{
		ImageID:     "8d8e8888-8ed8-88b8-88cb-888f8b8cf8fa",
		ImageName:   "CTyunOS-test",
		RegionID:    "88f8888888dd88ec888888888d888d8b",
		CmkID:       "9b5e62ab-015c-41d0-a7d6-f9db9d7fecXXX",
		Description: "Test CTyunOS",
		Labels: []*CtimageCopyImageLabelsRequest{
			{
				LabelKey:   "test-key",
				LabelValue: "test-value",
			},
		},
		ProjectID: "0",
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
