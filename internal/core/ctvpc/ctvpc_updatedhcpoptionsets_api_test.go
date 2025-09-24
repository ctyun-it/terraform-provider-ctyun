package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtvpcUpdatedhcpoptionsetsApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtvpcUpdatedhcpoptionsetsApi

	// 构造请求
	var name string = ""
	var description string = ""
	var domainName string = ""
	request := &CtvpcUpdatedhcpoptionsetsRequest{
		RegionID:         "",
		DhcpOptionSetsID: "",
		Name:             &name,
		Description:      &description,
		DomainName:       &domainName,
		DnsList:          []*string{},
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
