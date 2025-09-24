package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestDcs2CreateRedisTemplateApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.Dcs2CreateRedisTemplateApi

	// 构造请求
	request := &Dcs2CreateRedisTemplateRequest{
		RegionId:   "bb9fdb42056f11eda1610242ac110002",
		ProdInstId: "a0784f6c27cf0139e79e055f89f03f1d",
		Params: []*Dcs2CreateRedisTemplateParamsRequest{
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
