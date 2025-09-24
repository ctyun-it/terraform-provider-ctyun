package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestDcs2UnBindElasticIPApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.Dcs2UnBindElasticIPApi

	// 构造请求
	request := &Dcs2UnBindElasticIPRequest{
		RegionId:       "bb9fdb42056f11eda1610242ac110002",
		BindObjType:    4,
		ProdInstId:     "712041127816987648",
		ElasticIp:      "172.58.45.36",
		BindByPaasProd: 0,
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
