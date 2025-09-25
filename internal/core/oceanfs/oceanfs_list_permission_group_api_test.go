package oceanfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestOceanfsListPermissionGroupApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.OceanfsListPermissionGroupApi

	// 构造请求
	request := &OceanfsListPermissionGroupRequest{
		RegionID:            "参考[请求示例",
		PermissionGroupFuid: "参考[请求示例",
		PageSize:            10,
		PageNo:              1,
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
