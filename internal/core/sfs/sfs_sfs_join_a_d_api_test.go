package sfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestSfsSfsJoinADApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.SfsSfsJoinADApi

	// 构造请求
	var isAnonymousAcc bool = true
	request := &SfsSfsJoinADRequest{
		RegionID:       "参考[请求示例]",
		SfsUID:         "参考[请求示例]",
		IsAnonymousAcc: &isAnonymousAcc,
		Keytab:         "BQI******",
		KeytabMd5:      "7aef6360174a26d79da69ff944f80885",
		KeytabName:     "test",
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
