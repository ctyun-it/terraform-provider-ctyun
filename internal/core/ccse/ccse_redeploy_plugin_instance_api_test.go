package ccse

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCcseRedeployPluginInstanceApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CcseRedeployPluginInstanceApi

	// 构造请求
	request := &CcseRedeployPluginInstanceRequest{
		ClusterId:    "47281b02f87757478f20b1827c97cadf",
		PluginName:   "plugin_inst",
		RegionId:     "bb9fdb42056f11eda1610242ac110002",
		Namespace:    "kube-system",
		ChartName:    "demo-tpl",
		ChartVersion: "0.1.0",
		InstanceName: "demo-tpl-inst",
		Values:       "",
		ValuesJson:   "",
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
