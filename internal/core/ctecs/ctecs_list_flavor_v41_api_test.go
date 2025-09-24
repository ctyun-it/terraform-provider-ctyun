package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtecsListFlavorV41Api_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtecsListFlavorV41Api

	// 构造请求
	request := &CtecsListFlavorV41Request{
		RegionID:     "bb9fdb42056f11eda1610242ac110002",
		AzName:       "",
		FlavorType:   "GPU_N_PI7",
		FlavorName:   "pi7.4xlarge.4",
		FlavorCPU:    16,
		FlavorRAM:    32,
		FlavorArch:   "x86",
		FlavorSeries: "g",
		FlavorID:     "f02916cc-0445-be64-5e41-64019e95dc07",
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
