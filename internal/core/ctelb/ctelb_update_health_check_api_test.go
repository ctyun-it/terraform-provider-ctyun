package ctelb

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtelbUpdateHealthCheckApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtelbUpdateHealthCheckApi

	// 构造请求
	request := &CtelbUpdateHealthCheckRequest{
		ClientToken:       "79fa97e3-c48b-xxxx-9f46-6a13d8163678",
		RegionID:          "",
		ID:                "",
		HealthCheckID:     "hc-m2zb05f7s8",
		Name:              "acl11",
		Description:       "test",
		Timeout:           2,
		MaxRetry:          2,
		Interval:          5,
		HttpMethod:        "",
		HttpUrlPath:       "",
		HttpExpectedCodes: []string{},
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
