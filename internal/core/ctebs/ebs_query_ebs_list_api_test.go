package ctebs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestEbsQueryEbsListApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.EbsQueryEbsListApi

	// 构造请求
	var diskStatus string = "attached"
	var azName string = "az1"
	var projectID string = "0"
	var diskType string = "SAS"
	var diskMode string = "VBD"
	var multiAttach string = "true"
	var isSystemVolume string = "false"
	var isEncrypt string = "false"
	var queryContent string = "test"
	request := &EbsQueryEbsListRequest{
		RegionID:       "41f64827f25f468595ffa3a5deb5d15d",
		PageNo:         1,
		PageSize:       10,
		DiskStatus:     &diskStatus,
		AzName:         &azName,
		ProjectID:      &projectID,
		DiskType:       &diskType,
		DiskMode:       &diskMode,
		MultiAttach:    &multiAttach,
		IsSystemVolume: &isSystemVolume,
		IsEncrypt:      &isEncrypt,
		QueryContent:   &queryContent,
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
