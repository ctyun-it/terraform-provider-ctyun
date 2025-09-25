package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestDcs2ChangeAutoRenewStatusApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.Dcs2ChangeAutoRenewStatusApi

	// 构造请求
	request := &Dcs2ChangeAutoRenewStatusRequest{
		RegionId:            "bb9fdb42056f11eda1610242ac110002",
		AutoRenewCycleType:  "3",
		AutoRenewCycleCount: 2,
		AutoRenewStatus:     "1",
		Source:              "8",
		ProdInstIds:         []string{},
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
