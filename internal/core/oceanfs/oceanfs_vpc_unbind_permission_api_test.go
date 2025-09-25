package oceanfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestOceanfsVpcUnbindPermissionApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.OceanfsVpcUnbindPermissionApi

	// 构造请求
	request := &OceanfsVpcUnbindPermissionRequest{
		PermissionGroupFuid: "参考[请求示例]",
		RegionID:            "参考[请求示例]",
		SfsUID:              "参考[请求示例]",
		VpcID:               "参考[请求示例]",
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
