package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestDcs2EditRedisTemplateApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.Dcs2EditRedisTemplateApi

	// 构造请求
	var sysTemplate bool = false
	request := &Dcs2EditRedisTemplateRequest{
		RegionId: "bb9fdb42056f11eda1610242ac110002",
		Template: &Dcs2EditRedisTemplateTemplateRequest{
			Id:          "00ab67330baf46b292477bd6ba4f2bca",
			Name:        "CLASSIC-test-1",
			Description: "自研系列参数模板",
			CacheMode:   "CLASSIC",
			SysTemplate: &sysTemplate,
		},
		Params: []*Dcs2EditRedisTemplateParamsRequest{
			{
				ParamName:     "hash-max-ziplist-value",
				OriginalValue: "60",
				CurrentValue:  "66",
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
