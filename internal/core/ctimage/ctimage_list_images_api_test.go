package ctimage

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtimageListImagesApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	//credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	credential := core.CredentialFromEnv()
	apis := NewApis("https://ctimage-global.ctapi-test.ctyun.cn:21443", client)
	api := apis.CtimageListImagesApi

	//imageVisibilityCode := int32(-1)
	//pageNo := int32(0)
	pageSize := int32(0)
	// 构造请求
	request := &CtimageListImagesRequest{
		RegionID: "81f7728662dd11ec810800155d307d5b",
		//AzName:     "az1",
		//CwaiType:   "control",
		//FlavorName: "s7.small.1",
		//ImageScene:          "dev",
		//ImageStatus: "rejected",
		//ImageSubcategory:    "app",
		//ImageVisibilityCode: &imageVisibilityCode,
		//PageNo: &pageNo,
		PageSize: &pageSize,
		//ProjectID:           "0",
		//QueryContent:        "CTyunOS",
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
