package ctebm

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestEbmImageListApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.EbmImageListApi

	// 构造请求
	var imageType string = "private"
	var imageUUID string = "im-rgk4fo8wo3spbf9p8geyurs1ydpx"
	var osName string = "windows"
	var osVersion string = "server-2016"
	var osType string = "windows"
	request := &EbmImageListRequest{
		RegionID:   "81f7728662dd11ec810800155d307d5b",
		AzName:     "az2",
		DeviceType: "physical.t2.large",
		ImageType:  &imageType,
		ImageUUID:  &imageUUID,
		OsName:     &osName,
		OsVersion:  &osVersion,
		OsType:     &osType,
		PageNo:     2,
		PageSize:   2,
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
