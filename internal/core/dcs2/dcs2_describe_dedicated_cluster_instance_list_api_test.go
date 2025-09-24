package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestDcs2DescribeDedicatedClusterInstanceListApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.Dcs2DescribeDedicatedClusterInstanceListApi

	// 构造请求
	request := &Dcs2DescribeDedicatedClusterInstanceListRequest{
		RegionId:     "bb9fdb42056f11eda1610242ac110002",
		PageIndex:    1,
		PageSize:     10,
		InstanceName: "idcsync-9624lxg",
		ProdInstId:   "ee2e22896a11ee43e9d379d1598d75a1",
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
