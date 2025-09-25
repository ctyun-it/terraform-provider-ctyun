package ctelb

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtelbCreateHealthCheckApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtelbCreateHealthCheckApi

	// 构造请求
	request := &CtelbCreateHealthCheckRequest{
		ClientToken:       "79fa97e3-c48b-xxxx-9f46-6a13d8163678",
		RegionID:          "",
		Name:              "acl11",
		Description:       "",
		Protocol:          "",
		Timeout:           2,
		Interval:          5,
		MaxRetry:          2,
		HttpMethod:        "",
		HttpUrlPath:       "",
		HttpExpectedCodes: []string{},
		ProtocolPort:      0,
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
