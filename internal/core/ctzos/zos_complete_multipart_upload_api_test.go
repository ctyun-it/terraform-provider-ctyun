package ctzos

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestZosCompleteMultipartUploadApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.ZosCompleteMultipartUploadApi

	// 构造请求
	request := &ZosCompleteMultipartUploadRequest{
		Bucket:   "bucket1",
		Key:      "obj1",
		RegionID: "332232eb-63aa-465e-9028-52e5123866f0",
		MultipartUpload: &ZosCompleteMultipartUploadMultipartUploadRequest{
			Parts: []*ZosCompleteMultipartUploadMultipartUploadPartsRequest{
				{
					PartNumber: 1,
					ETag:       "5f363e0e58a95f06cbe9bbc662c5dfb6",
				},
			},
		},
		UploadID: "2~3BRpGICgM_5Lym31H53HhYl-XD6vjeH",
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
