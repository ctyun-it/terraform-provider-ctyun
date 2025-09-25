package ctgkafka

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtgkafkaChaosInjectApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtgkafkaChaosInjectApi

	// 构造请求
	request := &CtgkafkaChaosInjectRequest{
		RegionId:   "bb9fdb42056f11eda1610242ac110002",
		ProdInstId: "68eef42fd8d042bb960d3c3244d9243e",
		ActionCode: "node-shutdown",
		ActionParameter: &CtgkafkaChaosInjectActionParameterRequest{
			CpuPercent:   50,
			Duration:     60,
			NodeKillType: 0,
			EcsId:        "61e9c073-d4fd-8a3a-ecaf-14d4a5382c39",
			AzName:       "cn-xinan1-3A",
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
