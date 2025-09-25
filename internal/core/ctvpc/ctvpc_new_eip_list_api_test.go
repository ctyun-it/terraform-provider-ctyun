package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtvpcNewEipListApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtvpcNewEipListApi

	// 构造请求
	var projectID string = ""
	var status string = ""
	var ipType string = ""
	var eipType string = ""
	var ip string = ""
	request := &CtvpcNewEipListRequest{
		ClientToken: "",
		RegionID:    "",
		ProjectID:   &projectID,
		Page:        1,
		PageNo:      1,
		PageSize:    1,
		Ids:         []*string{},
		Status:      &status,
		IpType:      &ipType,
		EipType:     &eipType,
		Ip:          &ip,
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
