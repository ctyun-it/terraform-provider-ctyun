package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtecsListInstanceStatusV41Api_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtecsListInstanceStatusV41Api

	// 构造请求
	request := &CtecsListInstanceStatusV41Request{
		RegionID:       "bb9fdb42056f11eda1610242ac110002",
		AzName:         "cn-huadong1-jsnj1A-public-ctcloud",
		InstanceIDList: "73f321ea-62ff-11ec-a8bc-005056898fe0,88f888ea-88ff-88ec-a8bc-888888888fe8",
		PageNo:         1,
		PageSize:       10,
		ProjectID:      "0",
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
