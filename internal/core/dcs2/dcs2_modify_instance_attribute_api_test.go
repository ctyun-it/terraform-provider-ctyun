package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestDcs2ModifyInstanceAttributeApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.Dcs2ModifyInstanceAttributeApi

	// 构造请求
	var protectionStatus bool = true
	request := &Dcs2ModifyInstanceAttributeRequest{
		RegionId:         "bb9fdb42056f11eda1610242ac110002",
		ProdInstId:       "873209915236156416",
		Description:      "该实例用于XX服务测试",
		ProtectionStatus: &protectionStatus,
		MaintenanceTime:  "03:00-05:00",
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
