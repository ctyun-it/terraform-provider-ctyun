package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtecsInstanceDetachSfsV41Api_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtecsInstanceDetachSfsV41Api

	// 构造请求
	var forceDel bool = false
	request := &CtecsInstanceDetachSfsV41Request{
		RegionID:   "bb9fdb42056f11eda1610242ac110002",
		InstanceID: "b67b7f1f-095b-1249-b379-8dd5cc542a05",
		SysInfoList: []*CtecsInstanceDetachSfsV41SysInfoListRequest{
			{
				FileSysRoute: "55.243.4.20:/mnt/sfs_cap/e3aacef1e0be40559b38a5c2158aa62e_mqfb3ej839cdz7ub",
				MountPoint:   "/mnt/docs",
			},
		},
		ForceDel: &forceDel,
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
