package ccse

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCcseGetPublicImageListApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CcseGetPublicImageListApi

	// 构造请求
	request := &CcseGetPublicImageListRequest{
		RegionId:   "bb9fdb42056f11eda1610242ac110002",
		FlavorName: "c7.2xlarge.4",
		VmType:     "ecs",
		ProjectId:  "0",
		AzName:     "cn-xinan1-3A",
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
