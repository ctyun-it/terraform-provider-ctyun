package ctzos

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestZosDeleteObjectsApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.ZosDeleteObjectsApi

	// 构造请求
	var quiet bool = true
	request := &ZosDeleteObjectsRequest{
		Bucket:   "bucket1",
		RegionID: "332232eb-63aa-465e-9028-52e5123866f0",
		Delete: &ZosDeleteObjectsDeleteRequest{
			Objects: []*ZosDeleteObjectsDeleteObjectsRequest{
				{
					Key:       "obj1",
					VersionID: "Mp6FKCm7LkZvCA3iTWj6sdqv.4vLAFS",
				},
			},
			Quiet: &quiet,
		},
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
